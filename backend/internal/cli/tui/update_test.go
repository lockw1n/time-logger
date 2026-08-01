package tui

import (
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/outbox"
)

func fixedNow() time.Time {
	return time.Date(2026, 7, 8, 10, 0, 0, 0, time.UTC) // Wednesday of the fixture week
}

// loadedModel returns a model with the fixture week already loaded, so tests can
// drive interactions without threading the initial fetch.
func loadedModel(client api.Client) Model {
	m := New(context.Background(), client, Options{CompanyID: 1, Date: "2026-07-08", Now: fixedNow})
	m.grid = buildGrid(weekFixture())
	m.loaded = true
	m.loading = false
	return m
}

// key builds a KeyMsg from a name, mapping the special keys the update logic
// distinguishes and treating everything else as literal runes.
func key(s string) tea.KeyMsg {
	switch s {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

func send(m Model, msg tea.Msg) (Model, tea.Cmd) {
	tm, cmd := m.Update(msg)
	return tm.(Model), cmd
}

// sendKeys threads a sequence of key names through Update, returning the final
// model and the last command produced.
func sendKeys(m Model, names ...string) (Model, tea.Cmd) {
	var cmd tea.Cmd
	for _, n := range names {
		m, cmd = send(m, key(n))
	}
	return m, cmd
}

func TestNavigationClamps(t *testing.T) {
	m := loadedModel(&fakeClient{})

	// Two rows: down moves to row 1, a second down clamps there.
	m, _ = sendKeys(m, "down")
	if m.cursor.row != 1 {
		t.Fatalf("after down, row = %d, want 1", m.cursor.row)
	}
	m, _ = sendKeys(m, "down")
	if m.cursor.row != 1 {
		t.Errorf("row should clamp at 1, got %d", m.cursor.row)
	}

	// Right five times reaches Sat(5); two more clamp at Sun(6).
	m, _ = sendKeys(m, "right", "right", "right", "right", "right", "right", "right")
	if m.cursor.day != 6 {
		t.Errorf("day should clamp at 6, got %d", m.cursor.day)
	}

	// Left past the start clamps at 0.
	m, _ = sendKeys(m, "left", "left", "left", "left", "left", "left", "left", "left")
	if m.cursor.day != 0 {
		t.Errorf("day should clamp at 0, got %d", m.cursor.day)
	}

	// hjkl mirror the arrows.
	m.cursor = cursor{row: 0, day: 0}
	m, _ = sendKeys(m, "j", "l")
	if m.cursor.row != 1 || m.cursor.day != 1 {
		t.Errorf("hjkl move = %+v, want {1,1}", m.cursor)
	}
}

// TestEnterDispatch covers the plan's cell-edit semantics: empty→add, single→edit,
// multi→entry list.
func TestEnterDispatch(t *testing.T) {
	base := loadedModel(&fakeClient{activities: []api.Activity{{ID: 1, Name: "dev"}}})

	// APP-123 (row 0) Monday (day 0): single entry → inline edit.
	m, _ := send(withCursor(base, 0, 0), key("enter"))
	if m.mode != modeEdit {
		t.Errorf("single-entry cell: mode = %v, want modeEdit", m.mode)
	}
	if m.edit.entryID != 42 {
		t.Errorf("edit target = %d, want 42", m.edit.entryID)
	}

	// APP-123 (row 0) Wednesday (day 2): empty cell → add form, prefilled ticket.
	m, _ = send(withCursor(base, 0, 2), key("enter"))
	if m.mode != modeAdd {
		t.Errorf("empty cell: mode = %v, want modeAdd", m.mode)
	}
	if m.add.ticketCode() != "APP-123" {
		t.Errorf("add form ticket = %q, want APP-123 (prefilled)", m.add.ticketCode())
	}

	// APP-200 (row 1) Wednesday (day 2): two entries → entry list.
	m, _ = send(withCursor(base, 1, 2), key("enter"))
	if m.mode != modeEntryList {
		t.Errorf("multi-entry cell: mode = %v, want modeEntryList", m.mode)
	}
	if len(m.list.entries) != 2 {
		t.Errorf("entry list has %d entries, want 2", len(m.list.entries))
	}
}

func TestDeleteDispatch(t *testing.T) {
	base := loadedModel(&fakeClient{})

	// Single entry → confirm.
	m, _ := send(withCursor(base, 0, 0), key("d"))
	if m.mode != modeConfirmDelete || m.confirm.entryID != 42 {
		t.Errorf("single: mode=%v id=%d, want confirm/42", m.mode, m.confirm.entryID)
	}

	// Empty cell → no mode change, a status note instead.
	m, _ = send(withCursor(base, 0, 2), key("d"))
	if m.mode != modeView || m.status == "" {
		t.Errorf("empty: mode=%v status=%q, want view + note", m.mode, m.status)
	}

	// Multi entry → entry list.
	m, _ = send(withCursor(base, 1, 2), key("d"))
	if m.mode != modeEntryList {
		t.Errorf("multi: mode=%v, want entry list", m.mode)
	}
}

// TestEditSubmit drives a single-cell duration edit end to end: type a value,
// enter, and confirm the resulting command PUTs the new minutes.
func TestEditSubmit(t *testing.T) {
	fc := &fakeClient{}
	m := loadedModel(fc)
	m, _ = send(withCursor(m, 0, 0), key("enter")) // open edit on entry 42

	// Clear the prefilled "2h" and type "3h".
	m.edit.input = newTextInput("")
	m, _ = sendKeys(m, "3", "h")

	m, cmd := send(m, key("enter"))
	if m.mode != modeView || !m.loading {
		t.Fatalf("after submit: mode=%v loading=%v, want view+loading", m.mode, m.loading)
	}

	// Execute the command and feed its message back, as the runtime would.
	msg := cmd()
	mm, ok := msg.(mutationMsg)
	if !ok || mm.err != nil || mm.verb != "updated" {
		t.Fatalf("mutation msg = %+v (ok=%v)", mm, ok)
	}
	if len(fc.updated) != 1 || fc.updated[0].id != 42 || *fc.updated[0].req.DurationMinutes != 180 {
		t.Fatalf("update call = %+v, want id 42 / 180m", fc.updated)
	}

	// The success message triggers a reconciling refetch.
	m2, cmd := send(m, mm)
	if !m2.loading || cmd == nil {
		t.Errorf("mutation success should refetch: loading=%v cmd=%v", m2.loading, cmd)
	}
	if reload, ok := cmd().(timesheetMsg); !ok {
		t.Errorf("refetch cmd produced %T, want timesheetMsg", reload)
	}
}

// TestEditEmptyDeletes covers "empty input = delete after confirm".
func TestEditEmptyDeletes(t *testing.T) {
	m := loadedModel(&fakeClient{})
	m, _ = send(withCursor(m, 0, 0), key("enter"))
	m.edit.input = newTextInput("")
	m, _ = send(m, key("enter"))
	if m.mode != modeConfirmDelete || m.confirm.entryID != 42 {
		t.Errorf("empty edit submit: mode=%v id=%d, want confirm/42", m.mode, m.confirm.entryID)
	}
}

// TestConfirmDelete confirms a delete issues the DELETE command.
func TestConfirmDelete(t *testing.T) {
	fc := &fakeClient{}
	m := loadedModel(fc)
	m, _ = send(withCursor(m, 0, 0), key("d")) // confirm delete of 42
	m, cmd := send(m, key("y"))
	if m.mode != modeView || !m.loading {
		t.Fatalf("after y: mode=%v loading=%v", m.mode, m.loading)
	}
	if msg, ok := cmd().(mutationMsg); !ok || msg.verb != "deleted" {
		t.Fatalf("delete cmd msg = %+v", msg)
	}
	if len(fc.deleted) != 1 || fc.deleted[0] != 42 {
		t.Errorf("deleted = %v, want [42]", fc.deleted)
	}

	// Declining leaves the entry alone.
	fc2 := &fakeClient{}
	m2 := loadedModel(fc2)
	m2, _ = send(withCursor(m2, 0, 0), key("d"))
	m2, _ = send(m2, key("n"))
	if m2.mode != modeView || len(fc2.deleted) != 0 {
		t.Errorf("decline: mode=%v deleted=%v", m2.mode, fc2.deleted)
	}
}

// TestAddSubmit fills the add form and confirms it submits a create request with
// the cursor's day.
func TestAddSubmit(t *testing.T) {
	fc := &fakeClient{activities: []api.Activity{{ID: 7, Name: "dev"}}}
	m := loadedModel(fc)
	m.activities = fc.activities

	m, _ = send(withCursor(m, 0, 2), key("a")) // add form, cursor on Wed (2026-07-08)
	if m.mode != modeAdd {
		t.Fatalf("mode = %v, want modeAdd", m.mode)
	}
	// Blank form: type a ticket, cycle to the activity field is preselected at 0.
	m.add.ticket = newTextInput("APP-500")
	m.add.duration = newTextInput("1h")

	m, cmd := send(m, key("enter"))
	if m.mode != modeView || !m.loading {
		t.Fatalf("after add submit: mode=%v loading=%v", m.mode, m.loading)
	}
	if msg, ok := cmd().(mutationMsg); !ok || msg.verb != "added" || msg.err != nil {
		t.Fatalf("add cmd msg = %+v", msg)
	}
	if len(fc.created) != 1 {
		t.Fatalf("created = %d requests, want 1", len(fc.created))
	}
	got := fc.created[0]
	if got.TicketCode != "APP-500" || got.ActivityID != 7 || got.DurationMinutes != 60 || got.Date != "2026-07-08" {
		t.Errorf("create req = %+v", got)
	}
}

// TestAddSubmitOfflineQueues: with the backend unreachable, submitting the add
// form queues the create to the offline outbox (matching `tl add`) instead of
// dropping the entry, and reports it as queued without triggering a reload.
func TestAddSubmitOfflineQueues(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	fc := &fakeClient{activities: []api.Activity{{ID: 7, Name: "dev"}}, createErr: api.ErrUnreachable}
	m := loadedModel(fc)
	m.activities = fc.activities

	m, _ = send(withCursor(m, 0, 2), key("a"))
	m.add.ticket = newTextInput("APP-500")
	m.add.duration = newTextInput("1h")

	m, cmd := send(m, key("enter"))
	msg, ok := cmd().(mutationMsg)
	if !ok || !msg.queued || msg.err != nil {
		t.Fatalf("add cmd msg = %+v, want queued with no error", msg)
	}

	// Feeding the queued result back must not put the model into a reload (still
	// offline) — a reload would just fail and clobber the queued status.
	m2, _ := send(m, msg)
	if m2.loading {
		t.Error("queued add should not leave the model loading (no reload while offline)")
	}

	if n, _ := outbox.Count(); n != 1 {
		t.Fatalf("outbox has %d ops after offline add, want 1", n)
	}
}

// TestWeekChangeClearsGrid guards the stale-grid bug: after paging to another
// week the old grid must be dropped so a delete/edit keypress during the fetch
// can't act on the previous week's entries.
func TestWeekChangeClearsGrid(t *testing.T) {
	fc := &fakeClient{}
	m := loadedModel(fc)
	if len(m.grid.rows) == 0 {
		t.Fatal("fixture should start with rows")
	}

	m, _ = send(m, key("]")) // next week — fetch in flight
	if len(m.grid.rows) != 0 {
		t.Fatalf("grid not cleared on week change: %d rows", len(m.grid.rows))
	}

	// A delete keypress before the new week loads must be a no-op, not a stale
	// delete of the previous week's entry.
	m, _ = send(withCursor(m, 0, 0), key("d"))
	if m.mode != modeView || len(fc.deleted) != 0 {
		t.Errorf("delete on cleared grid: mode=%v deleted=%v", m.mode, fc.deleted)
	}
}

// TestCompanySwitchClearsGrid is the company-switch half of the same guard.
func TestCompanySwitchClearsGrid(t *testing.T) {
	fc := &fakeClient{}
	m := loadedModel(fc)
	m.companies = []api.Company{{ID: 1, Name: "Acme"}, {ID: 2, Name: "Globex"}}

	m, _ = send(m, key("c"))
	if len(m.grid.rows) != 0 {
		t.Fatalf("grid not cleared on company switch: %d rows", len(m.grid.rows))
	}
	m, _ = send(withCursor(m, 0, 0), key("enter"))
	if m.mode == modeEdit {
		t.Error("enter opened an edit against the previous company's grid")
	}
}

// TestPrefilledAddUsesRowActivity guards against the wrong-activity default: when
// the row's activity isn't in the fetched list, the prefilled add form must still
// submit under the row's activity id, not activities[0].
func TestPrefilledAddUsesRowActivity(t *testing.T) {
	fc := &fakeClient{activities: []api.Activity{{ID: 9, Name: "other"}}}
	m := loadedModel(fc)
	m.activities = fc.activities

	// APP-123 (row 0, activity id 1) on an empty Wednesday cell (day 2).
	m, _ = send(withCursor(m, 0, 2), key("enter"))
	if m.mode != modeAdd {
		t.Fatalf("mode = %v, want modeAdd", m.mode)
	}
	m.add.duration = newTextInput("1h")

	_, cmd := send(m, key("enter"))
	if _, ok := cmd().(mutationMsg); !ok {
		t.Fatal("expected a create command")
	}
	if len(fc.created) != 1 {
		t.Fatalf("created %d entries, want 1", len(fc.created))
	}
	if got := fc.created[0]; got.ActivityID != 1 || got.TicketCode != "APP-123" {
		t.Errorf("create req = %+v, want ticket APP-123 / activity 1", got)
	}
}

func TestWeekNavigation(t *testing.T) {
	m := loadedModel(&fakeClient{})

	m2, cmd := send(m, key("]"))
	if m2.weekStart != "2026-07-13" || !m2.loading || cmd == nil {
		t.Errorf("next week: start=%s loading=%v", m2.weekStart, m2.loading)
	}

	m3, _ := send(m, key("["))
	if m3.weekStart != "2026-06-29" {
		t.Errorf("prev week: start=%s, want 2026-06-29", m3.weekStart)
	}

	// From a different week, `t` returns to the week containing now (fixedNow).
	m4 := m2
	m4, _ = send(m4, key("t"))
	if m4.weekStart != "2026-07-06" {
		t.Errorf("today: start=%s, want 2026-07-06", m4.weekStart)
	}
}

func TestCompanySwitch(t *testing.T) {
	m := loadedModel(&fakeClient{})
	m.companies = []api.Company{{ID: 1, Name: "Acme"}, {ID: 2, Name: "Globex"}}
	m.companyID = 1

	m, cmd := send(m, key("c"))
	if m.companyID != 2 || !m.loading || cmd == nil {
		t.Errorf("switch: companyID=%d loading=%v", m.companyID, m.loading)
	}

	// Single company → no switch, just a note.
	m2 := loadedModel(&fakeClient{})
	m2.companies = []api.Company{{ID: 1, Name: "Acme"}}
	m2, _ = send(m2, key("c"))
	if m2.companyID != 1 || m2.status == "" {
		t.Errorf("single company: id=%d status=%q", m2.companyID, m2.status)
	}
}

func TestQuit(t *testing.T) {
	m := loadedModel(&fakeClient{})
	_, cmd := send(m, key("q"))
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("q should quit")
	}
	_, cmd = send(m, key("ctrl+c"))
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("ctrl+c should quit")
	}
}

// TestStaleTimesheetDropped ensures a late reply for another week/company is
// ignored so it can't overwrite the view the user has moved to.
func TestStaleTimesheetDropped(t *testing.T) {
	m := loadedModel(&fakeClient{})
	before := m.grid

	m2, _ := send(m, timesheetMsg{companyID: 1, start: "2026-01-01", ts: api.TimesheetResponse{}})
	if len(m2.grid.rows) != len(before.rows) {
		t.Errorf("stale reply overwrote the grid")
	}
}

// TestTimesheetError shows an API error lands in the footer without changing mode.
func TestTimesheetError(t *testing.T) {
	m := loadedModel(&fakeClient{})
	m.loading = true
	m2, _ := send(m, timesheetMsg{companyID: 1, start: m.weekStart, err: api.ErrUnreachable})
	if m2.loading || !m2.statusErr || m2.status == "" {
		t.Errorf("error handling: loading=%v statusErr=%v status=%q", m2.loading, m2.statusErr, m2.status)
	}
	if m2.mode != modeView {
		t.Errorf("error should not change mode, got %v", m2.mode)
	}
}

func TestWindowSizeAndTick(t *testing.T) {
	m := loadedModel(&fakeClient{})
	m2, _ := send(m, tea.WindowSizeMsg{Width: 120, Height: 40})
	if m2.width != 120 || m2.height != 40 {
		t.Errorf("size not recorded: %dx%d", m2.width, m2.height)
	}
	if _, cmd := send(m2, tickMsg(time.Now())); cmd == nil {
		t.Errorf("tick should reschedule itself")
	}
}

// withCursor returns a copy of m with the cursor placed at (row, day).
func withCursor(m Model, row, day int) Model {
	m.cursor = cursor{row: row, day: day}
	return m
}
