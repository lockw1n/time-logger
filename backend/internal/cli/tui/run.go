package tui

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// Run opens the interactive weekly timesheet and blocks until the user quits.
//
// The alt-screen buffer and terminal raw mode are owned by Bubble Tea, which
// restores them on normal exit, on a cancelled context, and on a panic (it
// recovers, restores the terminal, then re-raises). The deferred recover here is
// the belt-and-braces final net: by the time it runs the terminal is already
// restored, so it just turns a would-be stack dump into a returned error and a
// clean exit code — a wrecked terminal after a crash is an instant-delete bug.
func Run(ctx context.Context, client api.Client, opts Options) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("ui crashed: %v", r)
		}
	}()

	// WithContext ties the program's lifetime to ctx: a SIGINT/SIGTERM that
	// cancels ctx (wired in commands.Execute) tears the UI down cleanly.
	p := tea.NewProgram(New(ctx, client, opts), tea.WithAltScreen(), tea.WithContext(ctx))
	if _, runErr := p.Run(); runErr != nil {
		// A context cancellation (Ctrl-C / signal) is a normal quit, not an error.
		if errors.Is(runErr, tea.ErrProgramKilled) || errors.Is(runErr, context.Canceled) {
			return nil
		}
		return runErr
	}
	return nil
}
