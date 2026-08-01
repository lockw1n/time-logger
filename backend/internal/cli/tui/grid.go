package tui

import (
	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/render"
)

// gridRow is one aggregated timesheet line — a ticket+activity pair — split into
// seven per-day buckets Mon..Sun. It mirrors the phase-2 `tl week` aggregation:
// rows are keyed by (ticket, activity), a cell holds the entries logged that day,
// and perDay is their summed minutes.
type gridRow struct {
	ticketCode   string
	activityName string
	activityID   uint64
	priority     int
	perDay       [7]int
	entries      [7][]api.TimesheetEntry
	total        int
}

// grid is the fully-derived view model for one week's timesheet: the seven day
// dates, the aggregated rows in display order, and the totals footer. It is a
// pure function of the API response so the renderer and the update logic share a
// single, testable projection.
type grid struct {
	days    [7]string // YYYY-MM-DD, Monday..Sunday
	rows    []gridRow
	totals  [7]int
	overall int
}

// buildGrid projects a timesheet response onto the weekly grid. The week's seven
// days come from ts.Start (the server always anchors it to Monday); rows are
// sorted by ticket code then activity priority to match `tl week` and the web UI.
func buildGrid(ts api.TimesheetResponse) grid {
	days := weekDays(ts.Start)

	dayIndex := make(map[string]int, 7)
	for i, d := range days {
		dayIndex[d] = i
	}

	// Project already-sorted rows so the grid shares render's exact ordering
	// (ticket code, then activity priority) rather than re-implementing it.
	rows := make([]gridRow, 0, len(ts.Rows))
	for _, row := range render.SortRows(ts.Rows) {
		gr := gridRow{
			ticketCode:   row.Ticket.Code,
			activityName: row.Activity.Name,
			activityID:   row.Activity.ID,
			priority:     row.Activity.Priority,
			total:        row.TotalMinutes,
		}
		for day, minutes := range row.PerDayMinutes {
			if i, ok := dayIndex[day]; ok {
				gr.perDay[i] = minutes
			}
		}
		for _, e := range row.Entries {
			if i, ok := dayIndex[e.Date]; ok {
				gr.entries[i] = append(gr.entries[i], e)
			}
		}
		rows = append(rows, gr)
	}

	var g grid
	g.days = days
	g.rows = rows
	for i, day := range days {
		g.totals[i] = ts.Totals.PerDayMinutes[day]
	}
	g.overall = ts.Totals.OverallMinutes
	return g
}

// weekDays returns the seven YYYY-MM-DD dates Monday..Sunday for the week that
// contains start. It falls back to start repeated on an unparseable date so the
// renderer degrades to an empty week rather than panicking.
func weekDays(start string) [7]string {
	var out [7]string
	days, err := parse.WeekDays(start)
	if err != nil {
		for i := range out {
			out[i] = start
		}
		return out
	}
	copy(out[:], days)
	return out
}

// entriesAt returns the timesheet entries logged in the cell at (row, day). It is
// the basis for the cell-edit dispatch: zero entries means the cell is empty
// (enter → add), exactly one means a direct edit/delete target, and two or more
// means the cell aggregates several entries and must open the entry list first.
func (g grid) entriesAt(row, day int) []api.TimesheetEntry {
	if row < 0 || row >= len(g.rows) || day < 0 || day > 6 {
		return nil
	}
	return g.rows[row].entries[day]
}
