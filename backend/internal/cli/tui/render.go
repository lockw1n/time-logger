package tui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/render"
)

// Layout constants for the grid. gutter is the blank space between columns;
// activityFloor is the narrowest an activity name is truncated to before the grid
// gives up on shrinking and the view falls back to the too-narrow notice.
const (
	gutter        = 2
	activityFloor = 6
	ellipsis      = "…"
)

var dayNames = [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// dayHeaders reports which of the seven columns is "today" (from date), for the
// emphasized header; -1 when today is outside the displayed week.
func todayColumn(g grid, today string) int {
	for i, d := range g.days {
		if d == today {
			return i
		}
	}
	return -1
}

// renderGrid draws the weekly grid as plain, column-aligned text. It is the
// golden-tested core: no lipgloss styling is applied here (the decoration lives
// in View), so its output is deterministic across environments. The selected
// cell is framed in [brackets] and the today column header is marked with a
// leading '*', both of which survive into a plain terminal or a pipe. width
// drives activity-name truncation so the table fits the terminal.
func renderGrid(g grid, cur cursor, todayCol, width int) string {
	if len(g.rows) == 0 {
		return fmt.Sprintf("no entries for the week of %s\n", g.days[0])
	}

	// Column widths. Ticket and total are short and never truncated; the activity
	// column absorbs any width pressure.
	ticketW := len("total")
	activityW := len("activity")
	totalW := len("total")
	dayW := make([]int, 7)
	for i := range dayW {
		dayW[i] = len(dayNames[i])
	}

	for _, row := range g.rows {
		ticketW = max(ticketW, dispWidth(row.ticketCode))
		activityW = max(activityW, dispWidth(row.activityName))
		totalW = max(totalW, dispWidth(parse.FormatMinutes(row.total)))
		for i, mins := range row.perDay {
			dayW[i] = max(dayW[i], dispWidth(render.Cell(mins)))
		}
	}
	for i, mins := range g.totals {
		dayW[i] = max(dayW[i], dispWidth(render.Cell(mins)))
	}
	totalW = max(totalW, dispWidth(parse.FormatMinutes(g.overall)))

	activityW = fitActivity(activityW, ticketW, totalW, dayW, width)

	var b strings.Builder

	// Header row.
	b.WriteString(pad("ticket", ticketW))
	b.WriteString(strings.Repeat(" ", gutter))
	b.WriteString(pad("activity", activityW))
	for i := range dayW {
		b.WriteString(strings.Repeat(" ", gutter))
		b.WriteString(frameHeader(dayNames[i], dayW[i], i == todayCol))
	}
	b.WriteString(strings.Repeat(" ", gutter))
	b.WriteString(pad("total", totalW))
	b.WriteByte('\n')

	// Data rows.
	for r, row := range g.rows {
		b.WriteString(pad(row.ticketCode, ticketW))
		b.WriteString(strings.Repeat(" ", gutter))
		b.WriteString(pad(truncate(row.activityName, activityW), activityW))
		for i := range dayW {
			b.WriteString(strings.Repeat(" ", gutter))
			selected := r == cur.row && i == cur.day
			b.WriteString(frameCell(render.Cell(row.perDay[i]), dayW[i], selected))
		}
		b.WriteString(strings.Repeat(" ", gutter))
		b.WriteString(pad(parse.FormatMinutes(row.total), totalW))
		b.WriteByte('\n')
	}

	// Totals footer.
	b.WriteString(pad("total", ticketW))
	b.WriteString(strings.Repeat(" ", gutter))
	b.WriteString(pad("", activityW))
	for i := range dayW {
		b.WriteString(strings.Repeat(" ", gutter))
		b.WriteString(frameCell(render.Cell(g.totals[i]), dayW[i], false))
	}
	b.WriteString(strings.Repeat(" ", gutter))
	b.WriteString(pad(parse.FormatMinutes(g.overall), totalW))
	b.WriteByte('\n')

	return b.String()
}

// fitActivity shrinks the activity column so the whole table fits width, down to
// activityFloor. It returns the original width when width is unknown (0) or the
// table already fits — the grid never widens a column to fill space.
func fitActivity(activityW, ticketW, totalW int, dayW []int, width int) int {
	if width <= 0 {
		return activityW
	}
	fixed := ticketW + gutter + totalW
	for _, w := range dayW {
		fixed += gutter + w + 2 // +2 for each day cell's frame
	}
	fixed += gutter // gutter before total
	// available for the activity column = width - everything else.
	avail := max(width-fixed, activityFloor)
	if avail < activityW {
		return avail
	}
	return activityW
}

// frameCell renders a day cell padded to w and framed: [value] when selected,
// otherwise a space on each side so every column aligns regardless of selection.
func frameCell(value string, w int, selected bool) string {
	padded := pad(value, w)
	if selected {
		return "[" + padded + "]"
	}
	return " " + padded + " "
}

// frameHeader renders a day-name header cell, matching frameCell's width. Today's
// column is marked with a leading '*' so the emphasis survives a colorless render.
func frameHeader(name string, w int, today bool) string {
	padded := pad(name, w)
	if today {
		return "*" + padded + " "
	}
	return " " + padded + " "
}

// dispWidth is the terminal display width of s in columns — not its byte length.
// It matters because the em-dash and ellipsis the grid uses are multi-byte runes
// that occupy a single column; padding by byte length would misalign every row
// containing an empty-day cell.
func dispWidth(s string) int { return runewidth.StringWidth(s) }

// pad left-aligns s in a field of w display columns (no truncation — callers
// truncate first where needed).
func pad(s string, w int) string {
	gap := w - dispWidth(s)
	if gap <= 0 {
		return s
	}
	return s + strings.Repeat(" ", gap)
}

// truncate shortens s to w display columns, appending an ellipsis when it would
// overflow (runewidth.Truncate accounts for the ellipsis's own width).
func truncate(s string, w int) string {
	return runewidth.Truncate(s, w, ellipsis)
}
