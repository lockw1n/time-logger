package tui

import "github.com/charmbracelet/lipgloss"

// Styles are intentionally minimal: the plan scopes themes and color config out,
// so this is one restrained palette used only for the header/footer chrome. The
// grid itself is rendered as plain text (the selected cell is framed in brackets,
// today's column marked with '*') so its golden tests stay deterministic without
// depending on lipgloss's terminal-profile detection.
var (
	headerStyle = lipgloss.NewStyle().Bold(true)
	dimStyle    = lipgloss.NewStyle().Faint(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green

	// cursorStyle is the block caret drawn inside a focused text field.
	cursorStyle = lipgloss.NewStyle().Reverse(true)
)
