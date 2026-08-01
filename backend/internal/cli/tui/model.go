// Package tui implements `tl ui`: a full-screen, interactive weekly timesheet
// built on Bubble Tea. It is strictly a view/edit layer over the existing API
// client — no new business logic. Writes go through the same seams the one-shot
// commands use (entrylog.Submit for create/merge, the api.Client for update and
// delete), so the TUI can never violate a server constraint the CLI already
// respects. The model follows the standard Elm loop: Update is a pure function of
// (Model, Msg), and every network call is an async tea.Cmd whose result comes
// back as a message, so the UI never blocks on I/O.
package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/entrylog"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

// minWidth is the narrowest terminal the grid is drawn in. Below it the columns
// can't be laid out sensibly, so the View shows a "too narrow" notice instead of
// a broken table.
const minWidth = 100

type mode int

const (
	modeView mode = iota
	modeEdit
	modeAdd
	modeConfirmDelete
	modeEntryList
)

// cursor addresses a grid cell: a row (ticket+activity) and a day column 0..6.
type cursor struct {
	row int
	day int
}

// Model is the complete TUI state. It is a plain value — no pointers into shared
// mutable state — so Update can be tested by feeding key messages and asserting
// the returned model, with no terminal involved.
type Model struct {
	ctx    context.Context
	client api.Client
	now    func() time.Time

	companyID   uint64
	companyName string
	companies   []api.Company
	activities  []api.Activity // for companyID; feeds the add form

	weekStart string // Monday YYYY-MM-DD
	weekEnd   string // Sunday YYYY-MM-DD
	grid      grid
	loaded    bool

	cursor cursor
	mode   mode

	edit    editState
	add     addForm
	confirm confirmState
	list    listState

	// header caches the timer/outbox status shown in the header. It is refreshed
	// by a command on each tick rather than read from disk on every render, so a
	// burst of keystrokes doesn't trigger a burst of file I/O.
	header headerInfo

	width  int
	height int

	loading   bool // an API call is in flight
	status    string
	statusErr bool

	quitting   bool
	panicOnKey bool
}

// editState is the inline duration editor for a single entry. An empty value
// submitted means "delete this entry" (confirmed first).
type editState struct {
	entryID  uint64
	ticket   string
	activity string
	day      int
	input    textInput
}

// confirmState is a yes/no delete confirmation for one entry.
type confirmState struct {
	entryID uint64
	label   string
}

// headerInfo is the cached, read-only timer/outbox status drawn in the header's
// second line. The elapsed time is derived from startedAt at render time so it
// stays live between the once-a-second refreshes without re-reading the file.
type headerInfo struct {
	timerRunning bool
	ticket       string
	activity     string
	startedAt    time.Time
	outboxN      int
}

// listState is the mini entry-list shown for a cell that aggregates 2+ entries,
// so the user picks which one to edit or delete instead of blind-editing.
type listState struct {
	entries  []api.TimesheetEntry
	idx      int
	ticket   string
	activity string
	day      int
}

// Options configure a session. Date is any day in the week to open at (default
// today); Now is injected so week math and the elapsed-timer display are testable.
// PanicOnKey is the temporary test hook the plan calls for: with it set the model
// panics on the next keypress, so a manual tester can confirm the terminal is
// restored after a crash (see `tl ui --panic-test`).
type Options struct {
	CompanyID  uint64
	Date       string
	Now        func() time.Time
	PanicOnKey bool
}

// New builds the initial model for a session. It resolves the opening week from
// Options.Date but does no I/O — the first fetches are issued by Init.
func New(ctx context.Context, client api.Client, opts Options) Model {
	now := opts.Now
	if now == nil {
		now = time.Now
	}

	date := opts.Date
	if date == "" {
		date = now().Format(parse.DateLayout)
	}
	monday, sunday, err := parse.WeekBounds(date)
	if err != nil {
		monday, sunday, _ = parse.WeekBounds(now().Format(parse.DateLayout))
	}

	return Model{
		ctx:        ctx,
		client:     client,
		now:        now,
		companyID:  opts.CompanyID,
		weekStart:  monday,
		weekEnd:    sunday,
		loading:    true,
		panicOnKey: opts.PanicOnKey,
	}
}

// Init kicks off the first loads: the week, the company list (for the header and
// the `c` switcher), the activities (for the add form), and the header tick.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchTimesheet(m.ctx, m.client, m.companyID, m.weekStart, m.weekEnd),
		fetchCompanies(m.ctx, m.client),
		fetchActivities(m.ctx, m.client, m.companyID),
		refreshHeader(),
		tick(),
	)
}

// Update is the pure heart of the loop. It never performs I/O directly: state
// transitions are computed here and any resulting work is returned as a tea.Cmd.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		// Re-arm the tick and refresh the (cached) header status off the render path.
		return m, tea.Batch(tick(), refreshHeader())

	case headerMsg:
		m.header = headerInfo(msg)
		return m, nil

	case timesheetMsg:
		return m.onTimesheet(msg), nil

	case companiesMsg:
		return m.onCompanies(msg), nil

	case activitiesMsg:
		return m.onActivities(msg), nil

	case mutationMsg:
		return m.onMutation(msg)

	case tea.KeyMsg:
		return m.onKey(msg)
	}

	return m, nil
}

func (m Model) onTimesheet(msg timesheetMsg) Model {
	// Drop a stale reply: the user paged to another week or switched company
	// while this request was in flight.
	if msg.companyID != m.companyID || msg.start != m.weekStart {
		return m
	}
	m.loading = false
	if msg.err != nil {
		m.setError(msg.err)
		return m
	}
	m.grid = buildGrid(msg.ts)
	m.loaded = true
	m.clampCursor()
	return m
}

func (m Model) onCompanies(msg companiesMsg) Model {
	if msg.err != nil {
		// A failed company list is non-fatal — the header just shows the id.
		m.setError(msg.err)
		return m
	}
	m.companies = msg.companies
	m.companyName = companyName(msg.companies, m.companyID)
	return m
}

func (m Model) onActivities(msg activitiesMsg) Model {
	if msg.companyID != m.companyID {
		return m
	}
	if msg.err != nil {
		m.setError(msg.err)
		return m
	}
	m.activities = msg.activities
	return m
}

func (m Model) onMutation(msg mutationMsg) (tea.Model, tea.Cmd) {
	m.loading = false
	if msg.err != nil {
		m.setError(msg.err)
		return m, nil
	}
	if msg.queued {
		// Offline: the create is safely queued. A reload would just fail (still
		// offline) and clobber this message, so only refresh the header's outbox
		// count and leave the grid as-is; the entry appears after the next sync.
		m.setStatus(msg.verb)
		return m, refreshHeader()
	}
	// Reconcile against the server: refetch the week so totals and cell contents
	// reflect exactly what was committed.
	m.setStatus(msg.verb)
	return m.reload()
}

// reload issues a background refetch of the current week and marks the model
// loading. It is the single reconciliation path used after every mutation and by
// the manual refresh key.
func (m Model) reload() (Model, tea.Cmd) {
	m.loading = true
	return m, fetchTimesheet(m.ctx, m.client, m.companyID, m.weekStart, m.weekEnd)
}

// setError records err as the transient footer line, styled as an error. It never
// stops the loop — a failed API call is reported, not fatal.
func (m *Model) setError(err error) {
	m.status = err.Error()
	m.statusErr = true
}

func (m *Model) setStatus(s string) {
	m.status = s
	m.statusErr = false
}

// clampCursor keeps the cursor inside the current grid after a refetch shrinks
// the row set (e.g. an entry was deleted, or the week changed).
func (m *Model) clampCursor() {
	if len(m.grid.rows) == 0 {
		m.cursor = cursor{}
		return
	}
	if m.cursor.row >= len(m.grid.rows) {
		m.cursor.row = len(m.grid.rows) - 1
	}
	if m.cursor.row < 0 {
		m.cursor.row = 0
	}
	if m.cursor.day < 0 {
		m.cursor.day = 0
	}
	if m.cursor.day > 6 {
		m.cursor.day = 6
	}
}

// companyName finds a company's display name by id, falling back to "" so the
// header degrades to just the id when the list hasn't loaded.
func companyName(companies []api.Company, id uint64) string {
	for _, c := range companies {
		if c.ID == id {
			return c.Name
		}
	}
	return ""
}

// dayDate returns the YYYY-MM-DD for the cursor's day column, or today's date
// when the week hasn't loaded yet (so an add issued early still targets a day).
func (m Model) dayDate(day int) string {
	if m.loaded && day >= 0 && day <= 6 {
		return m.grid.days[day]
	}
	return m.now().Format(parse.DateLayout)
}

// buildRequest assembles the entry payload for a create/merge submission from the
// add form and the current company.
func (m Model) buildRequest(f addForm) (entrylog.Request, error) {
	minutes, err := parse.ParseDuration(f.duration.value)
	if err != nil {
		return entrylog.Request{}, err
	}
	activity := f.activities[f.activityIdx]
	return entrylog.Request{
		CompanyID:  m.companyID,
		TicketCode: f.ticketCode(),
		Activity:   activity,
		Date:       f.date,
		Minutes:    minutes,
		Comment:    f.commentPtr(),
	}, nil
}
