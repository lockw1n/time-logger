// Package render turns a timesheet response into the aligned text the `today`
// and `week` commands print. The functions are pure so their output can be
// asserted with golden strings.
package render

import (
	"sort"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// emptyCell is printed for a day with no logged time.
const emptyCell = "—"

// SortRows orders rows by ticket code, then by activity priority — matching the
// web UI. It sorts a copy so the caller's slice is left untouched. Exported so
// the TUI shares the exact same row ordering instead of re-implementing it.
func SortRows(rows []api.TimesheetRow) []api.TimesheetRow {
	out := make([]api.TimesheetRow, len(rows))
	copy(out, rows)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Ticket.Code != out[j].Ticket.Code {
			return out[i].Ticket.Code < out[j].Ticket.Code
		}
		return out[i].Activity.Priority < out[j].Activity.Priority
	})
	return out
}
