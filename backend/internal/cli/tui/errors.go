package tui

import "errors"

// User-facing validation errors surfaced in the footer status line. They mirror
// the guidance the one-shot commands give, so the TUI never fails silently.
var (
	errNoTicket     = errors.New("ticket code is required")
	errNoActivities = errors.New("no activities loaded — press r to retry, or run 'tl activities' online")
)
