package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

// onKey dispatches a keypress to the handler for the current mode. Ctrl-C always
// quits, from any mode, so the terminal is never trapped.
func (m Model) onKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		m.quitting = true
		return m, tea.Quit
	}
	if m.panicOnKey {
		panic("forced panic (test hook) — terminal should still be restored")
	}
	switch m.mode {
	case modeEdit:
		return m.onKeyEdit(msg)
	case modeAdd:
		return m.onKeyAdd(msg)
	case modeConfirmDelete:
		return m.onKeyConfirm(msg)
	case modeEntryList:
		return m.onKeyList(msg)
	default:
		return m.onKeyView(msg)
	}
}

// onKeyView handles navigation and the action keys in the main grid.
func (m Model) onKeyView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "up", "k":
		m.moveRow(-1)
	case "down", "j":
		m.moveRow(1)
	case "left", "h":
		m.moveDay(-1)
	case "right", "l":
		m.moveDay(1)
	case "[":
		return m.shiftWeek(-1)
	case "]":
		return m.shiftWeek(1)
	case "t":
		return m.gotoWeek(m.now().Format(parse.DateLayout))
	case "r":
		m.setStatus("refreshing…")
		return m.reload()
	case "c":
		return m.cycleCompany()
	case "a":
		return m.openAddBlank(), nil
	case "d":
		return m.deleteAction(), nil
	case "enter":
		return m.cellAction()
	}
	return m, nil
}

func (m *Model) moveRow(delta int) {
	if len(m.grid.rows) == 0 {
		return
	}
	m.cursor.row += delta
	m.clampCursor()
}

func (m *Model) moveDay(delta int) {
	m.cursor.day += delta
	m.clampCursor()
}

// shiftWeek pages the view delta weeks forward/back and refetches.
func (m Model) shiftWeek(delta int) (tea.Model, tea.Cmd) {
	d, err := time.Parse(parse.DateLayout, m.weekStart)
	if err != nil {
		return m, nil
	}
	return m.gotoWeek(d.AddDate(0, 0, delta*7).Format(parse.DateLayout))
}

// gotoWeek jumps to the week containing date and refetches it.
func (m Model) gotoWeek(date string) (tea.Model, tea.Cmd) {
	monday, sunday, err := parse.WeekBounds(date)
	if err != nil {
		return m, nil
	}
	m.weekStart, m.weekEnd = monday, sunday
	// Drop the old week's grid immediately: it's hidden behind the "loading"
	// notice, but leaving it in place would let an edit/delete keypress during the
	// fetch act on the previous week's entries (they have real, unrelated IDs).
	m.grid = grid{}
	m.loaded = false
	m.cursor.row = 0
	m.loading = true
	return m, fetchTimesheet(m.ctx, m.client, m.companyID, m.weekStart, m.weekEnd)
}

// cycleCompany switches to the next company and refetches both its week and its
// activities. A single-company consultant gets a status note instead.
func (m Model) cycleCompany() (tea.Model, tea.Cmd) {
	if len(m.companies) < 2 {
		m.setStatus("only one company")
		return m, nil
	}
	idx := 0
	for i, c := range m.companies {
		if c.ID == m.companyID {
			idx = i
			break
		}
	}
	next := m.companies[(idx+1)%len(m.companies)]
	m.companyID = next.ID
	m.companyName = next.Name
	m.activities = nil
	// Drop the old company's grid so a keypress during the fetch can't edit or
	// delete an entry belonging to the company we just switched away from.
	m.grid = grid{}
	m.loaded = false
	m.cursor.row = 0
	m.loading = true
	m.setStatus(fmt.Sprintf("company: %s", next.Name))
	return m, tea.Batch(
		fetchTimesheet(m.ctx, m.client, m.companyID, m.weekStart, m.weekEnd),
		fetchActivities(m.ctx, m.client, m.companyID),
	)
}

// cellAction runs the enter-key dispatch described in the plan: zero entries →
// add (prefilled from the row), exactly one → inline edit, two or more → the
// entry list so the user disambiguates before editing.
func (m Model) cellAction() (tea.Model, tea.Cmd) {
	if len(m.grid.rows) == 0 {
		return m.openAddBlank(), nil
	}
	entries := m.grid.entriesAt(m.cursor.row, m.cursor.day)
	switch len(entries) {
	case 0:
		return m.openAddPrefilled(), nil
	case 1:
		row := m.grid.rows[m.cursor.row]
		return m.openEdit(entries[0], row.ticketCode, row.activityName, m.cursor.day), nil
	default:
		return m.openList(entries), nil
	}
}

// deleteAction handles `d` on a cell: nothing on an empty cell, a confirm on a
// single entry, or the entry list when the cell aggregates several.
func (m Model) deleteAction() Model {
	if len(m.grid.rows) == 0 {
		return m
	}
	entries := m.grid.entriesAt(m.cursor.row, m.cursor.day)
	row := m.grid.rows[m.cursor.row]
	switch len(entries) {
	case 0:
		m.setStatus(fmt.Sprintf("nothing to delete on %s", m.dayDate(m.cursor.day)))
		return m
	case 1:
		return m.openConfirm(entries[0].ID, row.ticketCode, entries[0].DurationMinutes, m.cursor.day)
	default:
		return m.openList(entries)
	}
}

// openAddBlank opens an empty add form for the cursor's day. It needs the
// activities list loaded so the user can pick one.
func (m Model) openAddBlank() Model {
	m.mode = modeAdd
	m.add = newAddForm(m.dayDate(m.cursor.day), m.activities, "", nil)
	m.clearStatus()
	if len(m.activities) == 0 {
		m.setError(errNoActivities)
	}
	return m
}

// openAddPrefilled opens the add form seeded from the current row (its ticket and
// activity), leaving only the duration to type. The row's activity is passed
// explicitly so the form selects it even if the fetched activities list is empty
// or lacks it — no errNoActivities guard is needed here.
func (m Model) openAddPrefilled() Model {
	row := m.grid.rows[m.cursor.row]
	prefActivity := &api.Activity{ID: row.activityID, Name: row.activityName}
	m.mode = modeAdd
	m.add = newAddForm(m.dayDate(m.cursor.day), m.activities, row.ticketCode, prefActivity)
	m.clearStatus()
	return m
}

// openEdit starts the inline duration editor for one entry, prefilled with its
// current duration in editable form.
func (m Model) openEdit(entry api.TimesheetEntry, ticket, activity string, day int) Model {
	m.mode = modeEdit
	m.edit = editState{
		entryID:  entry.ID,
		ticket:   ticket,
		activity: activity,
		day:      day,
		input:    newTextInput(parse.FormatMinutes(entry.DurationMinutes)),
	}
	m.clearStatus()
	return m
}

// openConfirm opens a yes/no delete confirmation for one entry.
func (m Model) openConfirm(entryID uint64, ticket string, minutes, day int) Model {
	m.mode = modeConfirmDelete
	m.confirm = confirmState{
		entryID: entryID,
		label: fmt.Sprintf("delete %s on %s (%s)?",
			ticket, m.dayDate(day), parse.FormatMinutes(minutes)),
	}
	m.clearStatus()
	return m
}

// openList opens the mini entry-list for a multi-entry cell.
func (m Model) openList(entries []api.TimesheetEntry) Model {
	row := m.grid.rows[m.cursor.row]
	m.mode = modeEntryList
	m.list = listState{
		entries:  entries,
		idx:      0,
		ticket:   row.ticketCode,
		activity: row.activityName,
		day:      m.cursor.day,
	}
	m.clearStatus()
	return m
}

func (m *Model) clearStatus() {
	m.status = ""
	m.statusErr = false
}

// onKeyEdit handles the inline duration editor. Enter with an empty field turns
// into a delete confirmation (per the plan); a valid duration is submitted.
func (m Model) onKeyEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = modeView
		return m, nil
	case tea.KeyEnter:
		if m.edit.input.trimmedEmpty() {
			// Empty value means "delete this entry" — route straight to the
			// confirmation (which sets modeConfirmDelete itself).
			return m.openConfirmFromEdit(), nil
		}
		minutes, err := parse.ParseDuration(m.edit.input.value)
		if err != nil {
			m.setError(err)
			return m, nil
		}
		id := m.edit.entryID
		m.mode = modeView
		m.loading = true
		m.setStatus("saving…")
		return m, updateDuration(m.ctx, m.client, id, minutes)
	default:
		m.edit.input.update(msg)
		return m, nil
	}
}

// openConfirmFromEdit routes an empty edit submission to the delete confirmation,
// carrying over the entry being edited.
func (m Model) openConfirmFromEdit() Model {
	m.mode = modeConfirmDelete
	m.confirm = confirmState{
		entryID: m.edit.entryID,
		label:   fmt.Sprintf("delete %s on %s?", m.edit.ticket, m.dayDate(m.edit.day)),
	}
	return m
}

// onKeyAdd drives the add form: tab/arrows move between fields, left/right cycle
// the activity, enter submits, esc cancels.
func (m Model) onKeyAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = modeView
		return m, nil
	case tea.KeyEnter:
		return m.submitAdd()
	case tea.KeyTab, tea.KeyDown:
		m.add.field = (m.add.field + 1) % addFieldCount
		return m, nil
	case tea.KeyShiftTab, tea.KeyUp:
		m.add.field = (m.add.field + addFieldCount - 1) % addFieldCount
		return m, nil
	}

	if m.add.field == fieldActivity {
		switch msg.String() {
		case "left":
			m.add.cycleActivity(-1)
		case "right":
			m.add.cycleActivity(1)
		}
		return m, nil
	}

	if in := m.add.focusedInput(); in != nil {
		in.update(msg)
	}
	return m, nil
}

// submitAdd validates the form and, if it holds, issues the create/merge command.
func (m Model) submitAdd() (tea.Model, tea.Cmd) {
	if m.add.ticketCode() == "" {
		m.setError(errNoTicket)
		return m, nil
	}
	if len(m.add.activities) == 0 {
		m.setError(errNoActivities)
		return m, nil
	}
	req, err := m.buildRequest(m.add)
	if err != nil {
		m.setError(err)
		return m, nil
	}
	m.mode = modeView
	m.loading = true
	m.setStatus("saving…")
	return m, createEntry(m.ctx, m.client, req)
}

// onKeyConfirm handles the delete confirmation.
func (m Model) onKeyConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		id := m.confirm.entryID
		m.mode = modeView
		m.loading = true
		m.setStatus("deleting…")
		return m, deleteEntry(m.ctx, m.client, id)
	case "n", "N", "esc", "q":
		m.mode = modeView
		m.setStatus("cancelled")
		return m, nil
	}
	return m, nil
}

// onKeyList handles the multi-entry cell list: move the selection, then edit or
// delete the chosen entry.
func (m Model) onKeyList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.list.idx > 0 {
			m.list.idx--
		}
	case "down", "j":
		if m.list.idx < len(m.list.entries)-1 {
			m.list.idx++
		}
	case "enter", "e":
		sel := m.list.entries[m.list.idx]
		return m.openEdit(sel, m.list.ticket, m.list.activity, m.list.day), nil
	case "d":
		sel := m.list.entries[m.list.idx]
		return m.openConfirm(sel.ID, m.list.ticket, sel.DurationMinutes, m.list.day), nil
	case "esc", "q":
		m.mode = modeView
	}
	return m, nil
}
