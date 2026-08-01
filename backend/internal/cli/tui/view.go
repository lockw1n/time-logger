package tui

import (
	"fmt"
	"strings"

	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

// viewHints is the always-visible keybinding legend for the grid.
const viewHints = "←→↑↓ move · enter edit · a add · d delete · [ ] week · t today · c company · r refresh · q quit"

// View renders the whole screen: a two-line header, the grid, and a footer whose
// bottom block changes with the mode. It reads the timer and outbox files each
// render (they are cheap and read-only) so the header always reflects reality.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.width > 0 && m.width < minWidth {
		return fmt.Sprintf(
			"terminal too narrow: %d cols, need at least %d.\nresize the window, or press q to quit.\n",
			m.width, minWidth)
	}

	var b strings.Builder
	b.WriteString(m.headerView())
	b.WriteString("\n\n")
	b.WriteString(m.gridView())
	b.WriteString("\n")
	b.WriteString(m.footerView())
	return b.String()
}

func (m Model) headerView() string {
	company := fmt.Sprintf("company: %d", m.companyID)
	if m.companyName != "" {
		company = fmt.Sprintf("company: %s (%d)", m.companyName, m.companyID)
	}
	title := headerStyle.Render(fmt.Sprintf("time-logger — week %s … %s", m.weekStart, m.weekEnd))

	var b strings.Builder
	fmt.Fprintf(&b, "%s    %s\n", title, company)
	b.WriteString(m.statusBarView())
	return b.String()
}

// statusBarView is the header's second line: the running timer (if any) and the
// pending outbox count (if any). It reads from the cached header snapshot
// (refreshed once a second by refreshHeader), not the filesystem, so it costs
// nothing to redraw on every keystroke. Elapsed time is derived from the cached
// start so it still advances live between refreshes.
func (m Model) statusBarView() string {
	var parts []string

	if m.header.timerRunning {
		elapsed := timer.ElapsedMinutes(m.now().Sub(m.header.startedAt))
		parts = append(parts, fmt.Sprintf("⏱ %s (%s) %s running",
			m.header.ticket, m.header.activity, parse.FormatMinutes(elapsed)))
	}

	if m.header.outboxN > 0 {
		parts = append(parts, fmt.Sprintf("outbox: %d pending", m.header.outboxN))
	}

	if len(parts) == 0 {
		return dimStyle.Render("no timer running")
	}
	return dimStyle.Render(strings.Join(parts, "    "))
}

func (m Model) gridView() string {
	if !m.loaded {
		if m.loading {
			return dimStyle.Render("loading…")
		}
		return dimStyle.Render("no data — press r to refresh")
	}
	today := m.now().Format(parse.DateLayout)
	return renderGrid(m.grid, m.cursor, todayColumn(m.grid, today), m.gridWidth())
}

// gridWidth is the width budget passed to the renderer, defaulting generously
// when the terminal size isn't known yet (e.g. before the first WindowSizeMsg).
func (m Model) gridWidth() int {
	if m.width > 0 {
		return m.width
	}
	return 0 // 0 => renderer keeps full column widths
}

// footerView renders the mode-specific block followed by the transient status
// line, so a message from an API command or a validation error is always shown.
func (m Model) footerView() string {
	var b strings.Builder

	switch m.mode {
	case modeEdit:
		b.WriteString(m.editView())
	case modeAdd:
		b.WriteString(m.addView())
	case modeConfirmDelete:
		b.WriteString(m.confirmView())
	case modeEntryList:
		b.WriteString(m.listView())
	default:
		b.WriteString(dimStyle.Render(viewHints))
	}
	b.WriteByte('\n')

	if m.loading {
		b.WriteString(dimStyle.Render("· working…"))
		b.WriteByte('\n')
	}
	if m.status != "" {
		style := statusStyle
		if m.statusErr {
			style = errorStyle
		}
		b.WriteString(style.Render(m.status))
		b.WriteByte('\n')
	}
	return b.String()
}

func (m Model) editView() string {
	e := m.edit
	return fmt.Sprintf("edit %s (%s) on %s — duration: %s\n%s",
		e.ticket, e.activity, m.dayDate(e.day), e.input.render(true),
		dimStyle.Render("enter save · empty+enter delete · esc cancel"))
}

func (m Model) addView() string {
	f := m.add
	var b strings.Builder
	fmt.Fprintf(&b, "add entry — %s\n", f.date)
	b.WriteString(formLine("ticket", f.ticket.render(f.field == fieldTicket), f.field == fieldTicket))
	b.WriteString(formLine("activity", activityDisplay(f), f.field == fieldActivity))
	b.WriteString(formLine("duration", f.duration.render(f.field == fieldDuration), f.field == fieldDuration))
	b.WriteString(formLine("comment", f.comment.render(f.field == fieldComment), f.field == fieldComment))
	b.WriteString(dimStyle.Render("tab/↑↓ move · ←→ change activity · enter save · esc cancel"))
	return b.String()
}

// activityDisplay renders the activity selector, showing the cycle hint when it
// holds the focus and a fallback when no activities are loaded.
func activityDisplay(f addForm) string {
	name := f.currentActivityName()
	if name == "" {
		return "(none — press r in the grid to load)"
	}
	if f.field == fieldActivity {
		return fmt.Sprintf("%s  ←/→", name)
	}
	return name
}

// formLine renders one labeled form row, marking the focused one.
func formLine(label, value string, focused bool) string {
	marker := "  "
	if focused {
		marker = "▸ "
	}
	return fmt.Sprintf("%s%-9s %s\n", marker, label+":", value)
}

func (m Model) confirmView() string {
	return fmt.Sprintf("%s %s", m.confirm.label, dimStyle.Render("(y/n)"))
}

func (m Model) listView() string {
	l := m.list
	var b strings.Builder
	fmt.Fprintf(&b, "%s (%s) on %s — %d entries:\n", l.ticket, l.activity, m.dayDate(l.day), len(l.entries))
	for i, e := range l.entries {
		marker := "  "
		if i == l.idx {
			marker = "▸ "
		}
		comment := ""
		if c := commentText(e.Comment); c != "" {
			comment = "  " + c
		}
		fmt.Fprintf(&b, "%s#%d  %s%s\n", marker, e.ID, parse.FormatMinutes(e.DurationMinutes), comment)
	}
	b.WriteString(dimStyle.Render("↑↓ select · e/enter edit · d delete · esc back"))
	return b.String()
}

func commentText(c *string) string {
	if c == nil {
		return ""
	}
	return *c
}
