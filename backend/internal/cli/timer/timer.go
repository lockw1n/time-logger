// Package timer persists the single live-timer state as a file so a timer
// started in one terminal is visible from any other and survives the shell that
// started it. There is no daemon: elapsed time is computed on read as
// now - StartedAt. At most one timer exists at a time, enforced by an exclusive
// create (hardlink publish) rather than a racy check-then-write.
package timer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/config"
)

const timerFile = "timer.json"

var (
	// ErrNotRunning means no timer file exists.
	ErrNotRunning = errors.New("no timer running")
	// ErrRunning means a timer is already active, so Start refused.
	ErrRunning = errors.New("a timer is already running")
)

// CorruptError wraps a timer.json that exists but can't be parsed, so callers
// can tell a genuine timer apart from a damaged file and suggest `tl cancel`.
type CorruptError struct {
	Path string
	Err  error
}

func (e *CorruptError) Error() string {
	return fmt.Sprintf("timer file %s is unreadable (%v) — run 'tl cancel' to clear it", e.Path, e.Err)
}

func (e *CorruptError) Unwrap() error { return e.Err }

// State is the persisted timer, matching timer.json on disk. The activity is
// stored by both id and name: the id (resolved at start) lets `tl stop` build
// the entry without a network lookup, which is what makes stopping work — and
// queueable — while the backend is unreachable. The name is kept for display.
type State struct {
	TicketCode   string    `json:"ticket_code"`
	ActivityID   uint64    `json:"activity_id,omitempty"`
	ActivityName string    `json:"activity_name"`
	CompanyID    uint64    `json:"company_id"`
	StartedAt    time.Time `json:"started_at"`
	Note         string    `json:"note,omitempty"`
}

// Elapsed reports how long the timer has run as of now. now is a parameter so
// callers (and tests) inject the clock rather than the store reading it.
func (s State) Elapsed(now time.Time) time.Duration {
	return now.Sub(s.StartedAt)
}

// ElapsedMinutes rounds a duration to the nearest minute, with a floor of 1 so a
// just-started timer still logs a minute rather than zero.
func ElapsedMinutes(d time.Duration) int {
	m := int(math.Round(d.Minutes()))
	if m < 1 {
		return 1
	}
	return m
}

// Path returns the timer file location within the tl config dir (TL_CONFIG_DIR aware).
func Path() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, timerFile), nil
}

// Start persists s as the active timer, failing with ErrRunning if one already
// exists. It writes the full content to a temp file first, then publishes it
// with a hardlink: os.Link fails if the target already exists, so it is both the
// exclusive lock — of two concurrent Starts exactly one wins, no check-then-write
// race — and an atomic full-content publish. A concurrent reader therefore never
// sees the timer file empty or half-written.
func Start(s State) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding %s: %w", timerFile, err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".timer-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return fmt.Errorf("setting mode on temp file: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	// Publish atomically; an existing target means a timer is already running.
	if err := os.Link(tmp.Name(), path); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return ErrRunning
		}
		return fmt.Errorf("creating %s: %w", timerFile, err)
	}
	return nil
}

// Load reads the active timer. It returns ErrNotRunning when none exists and a
// *CorruptError when the file is present but unparseable or missing key fields.
func Load() (State, error) {
	path, err := Path()
	if err != nil {
		return State{}, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return State{}, ErrNotRunning
	}
	if err != nil {
		return State{}, fmt.Errorf("reading %s: %w", timerFile, err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, &CorruptError{Path: path, Err: err}
	}
	if s.TicketCode == "" || s.StartedAt.IsZero() {
		return State{}, &CorruptError{Path: path, Err: errors.New("missing ticket_code or started_at")}
	}
	return s, nil
}

// Clear removes the active timer, if any. A missing file is not an error, so it
// is safe to call on an already-stopped timer.
func Clear() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("removing %s: %w", timerFile, err)
	}
	return nil
}
