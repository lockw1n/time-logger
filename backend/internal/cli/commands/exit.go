package commands

import "fmt"

// exitError lets a command choose the process exit code while owning its own
// user-facing output — main prints nothing extra for it. Two uses: a "normal"
// non-zero signal (status with no timer running) and cases where the command has
// already written a tailored message (stop keeping the timer after a failure).
type exitError struct{ code int }

func (e exitError) Error() string { return fmt.Sprintf("exit status %d", e.code) }

// ExitCode is the process exit code main should use. main detects this via an
// anonymous interface, so it stays decoupled from this concrete type.
func (e exitError) ExitCode() int { return e.code }

func exitWith(code int) error { return exitError{code: code} }
