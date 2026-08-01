package render

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

var weekdayHeaders = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// Week renders a Monday–Sunday grid for the week covered by the timesheet, with
// per-day columns, a per-row total, and a totals footer row.
func Week(ts api.TimesheetResponse) string {
	var b strings.Builder

	days, err := parse.WeekDays(ts.Start)
	if err != nil {
		// ts.Start always comes from the server in YYYY-MM-DD; treat a parse
		// failure as an empty week rather than panicking.
		fmt.Fprintf(&b, "no entries for the week of %s\n", ts.Start)
		return b.String()
	}

	if len(ts.Rows) == 0 {
		fmt.Fprintf(&b, "no entries for the week of %s\n", days[0])
		return b.String()
	}

	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	// Header.
	fmt.Fprintf(w, "ticket\tactivity\t%s\ttotal\n", strings.Join(weekdayHeaders, "\t"))

	for _, row := range SortRows(ts.Rows) {
		cells := make([]string, len(days))
		for i, day := range days {
			cells[i] = Cell(row.PerDayMinutes[day])
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			row.Ticket.Code,
			row.Activity.Name,
			strings.Join(cells, "\t"),
			parse.FormatMinutes(row.TotalMinutes),
		)
	}

	// Totals footer.
	totals := make([]string, len(days))
	for i, day := range days {
		totals[i] = Cell(ts.Totals.PerDayMinutes[day])
	}
	fmt.Fprintf(w, "total\t\t%s\t%s\n",
		strings.Join(totals, "\t"),
		parse.FormatMinutes(ts.Totals.OverallMinutes),
	)

	w.Flush()
	return b.String()
}

// Cell renders a per-day minute count, using the em-dash placeholder for zero.
// Exported so the TUI grid shares the exact same empty-day rendering.
func Cell(minutes int) string {
	if minutes <= 0 {
		return emptyCell
	}
	return parse.FormatMinutes(minutes)
}
