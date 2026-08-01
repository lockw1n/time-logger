package render

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

// Today renders a single-day timesheet: one row per ticket/activity with its
// duration, comment and entry id, followed by a total line. The entry id column
// is what `edit`/`delete` (phase 4) consume.
func Today(ts api.TimesheetResponse) string {
	var b strings.Builder

	if len(ts.Rows) == 0 {
		fmt.Fprintf(&b, "no entries for %s\n", ts.Start)
		return b.String()
	}

	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for _, row := range SortRows(ts.Rows) {
		id, comment := dayEntry(row, ts.Start)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			row.Ticket.Code,
			row.Activity.Name,
			parse.FormatMinutes(row.TotalMinutes),
			comment,
			id,
		)
	}
	w.Flush()

	fmt.Fprintf(&b, "total: %s\n", parse.FormatMinutes(ts.Totals.OverallMinutes))
	return b.String()
}

// dayEntry returns the entry id (as "#N") and comment for the given date within
// a row. With the ux_entry unique index there is at most one entry per row/day.
func dayEntry(row api.TimesheetRow, date string) (id, comment string) {
	for _, e := range row.Entries {
		if e.Date == date {
			return fmt.Sprintf("#%d", e.ID), api.Deref(e.Comment)
		}
	}
	return "", ""
}
