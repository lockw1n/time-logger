package commands

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/auth"
	"github.com/lockw1n/time-logger/internal/cli/outbox"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

// exitCode extracts the exit code a command asked for via the exitError seam.
func exitCode(t *testing.T, err error) int {
	t.Helper()
	var coder interface{ ExitCode() int }
	if !errors.As(err, &coder) {
		t.Fatalf("error %v does not carry an exit code", err)
	}
	return coder.ExitCode()
}

// activitiesOnlyMux serves the activity lookup that start/stop both perform.
func activitiesOnlyMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/activities", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, []map[string]any{
			{"id": 5, "company_id": 1, "name": "dev", "billable": true, "priority": 1},
		})
	})
	return mux
}

// runCmd executes tl with args against an empty stdin, returning stdout, stderr
// and the command error.
func runCmd(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	var out, errBuf bytes.Buffer
	cmd := NewRoot("test")
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), errBuf.String(), err
}

func TestStatusNotRunningExitsOne(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	_, errBuf, err := runCmd(t, "status")
	if got := exitCode(t, err); got != 1 {
		t.Fatalf("status exit = %d, want 1", got)
	}
	if !strings.Contains(errBuf, "no timer running") {
		t.Errorf("stderr = %q, want 'no timer running'", errBuf)
	}
}

func TestStatusRunningShowsElapsed(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustStart(t, timer.State{TicketCode: "APP-9", ActivityName: "dev", CompanyID: 1, StartedAt: time.Now().Add(-47 * time.Minute)})

	out, _, err := runCmd(t, "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !strings.Contains(out, "APP-9") || !strings.Contains(out, "elapsed") {
		t.Errorf("status output = %q", out)
	}
}

// The core preservation guarantee: if the log API fails, the timer file must
// survive so nothing is lost and the user can retry.
func TestStopPreservesTimerOnAPIError(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))
	mustStart(t, timer.State{TicketCode: "APP-123", ActivityName: "dev", CompanyID: 1, StartedAt: time.Now().Add(-30 * time.Minute)})

	mux := activitiesOnlyMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "boom"})
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	_, _, err := runCmd(t, "stop", "--base-url", ts.URL, "--yes")
	if got := exitCode(t, err); got != 1 {
		t.Fatalf("stop exit = %d, want 1", got)
	}
	if _, loadErr := timer.Load(); loadErr != nil {
		t.Fatalf("timer should be preserved after API failure, Load = %v", loadErr)
	}
}

// An unreachable backend must not lose the measured time: stop queues the entry
// to the outbox, drops the timer (its job is done — the minutes are in the
// payload), and exits 3 so scripts can tell "accepted but not confirmed" apart.
func TestStopQueuesEntryWhenUnreachable(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TL_CONFIG_DIR", dir)
	mustSaveToken(t, time.Now().Add(time.Hour))
	// The activity id is stored at start, so stop needs no network lookup and can
	// reach the create (which is what fails offline and gets queued).
	mustStart(t, timer.State{TicketCode: "APP-123", ActivityID: 5, ActivityName: "dev", CompanyID: 1, StartedAt: time.Now().Add(-10 * time.Minute)})

	// A server that is immediately closed → connection refused.
	ts := httptest.NewServer(activitiesOnlyMux())
	url := ts.URL
	ts.Close()

	out, _, err := runCmd(t, "stop", "--base-url", url, "--yes")
	if got := exitCode(t, err); got != 3 {
		t.Fatalf("stop exit = %d, want 3 (queued)", got)
	}
	if !strings.Contains(out, "queued") {
		t.Errorf("stdout = %q, want it to mention the entry was queued", out)
	}
	if _, loadErr := timer.Load(); !errors.Is(loadErr, timer.ErrNotRunning) {
		t.Fatalf("timer should be cleared after queueing, Load = %v", loadErr)
	}
	pending, perr := outbox.Pending()
	if perr != nil {
		t.Fatalf("reading outbox: %v", perr)
	}
	if len(pending) != 1 {
		t.Fatalf("outbox has %d ops, want 1", len(pending))
	}
	if got := pending[0].Op.Payload.TicketCode; got != "APP-123" {
		t.Errorf("queued op ticket = %q, want APP-123", got)
	}
}

// An expired token keeps the timer and exits 2 with the login-and-retry message.
func TestStopExpiredTokenKeepsTimerExitsTwo(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(-time.Hour)) // already expired
	mustStart(t, timer.State{TicketCode: "APP-1", ActivityName: "dev", CompanyID: 1, StartedAt: time.Now().Add(-10 * time.Minute)})

	// tokenFn fails before any request, so the base URL is never dialed.
	_, errBuf, err := runCmd(t, "stop", "--base-url", "http://127.0.0.1:0", "--yes")
	if got := exitCode(t, err); got != 2 {
		t.Fatalf("stop exit = %d, want 2", got)
	}
	if !strings.Contains(errBuf, "tl login") {
		t.Errorf("stderr = %q, want it to mention 'tl login'", errBuf)
	}
	if _, loadErr := timer.Load(); loadErr != nil {
		t.Fatalf("timer should be preserved on expired token, Load = %v", loadErr)
	}
}

func TestStartWritesTimerAndRefusesSecond(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))

	ts := httptest.NewServer(activitiesOnlyMux())
	defer ts.Close()

	if _, _, err := runCmd(t, "start", "APP-7", "-a", "dev", "--base-url", ts.URL, "--company", "1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	state, err := timer.Load()
	if err != nil {
		t.Fatalf("timer not written: %v", err)
	}
	if state.TicketCode != "APP-7" || state.ActivityName != "dev" {
		t.Errorf("stored state = %+v, want APP-7/dev", state)
	}

	// A second start while one runs is refused, naming the running ticket.
	_, _, err = runCmd(t, "start", "APP-8", "-a", "dev", "--base-url", ts.URL, "--company", "1")
	if err == nil {
		t.Fatal("second start should fail while a timer runs")
	}
	if !strings.Contains(err.Error(), "APP-7") || !strings.Contains(err.Error(), "running") {
		t.Errorf("error = %q, want it to name the running APP-7 timer", err.Error())
	}
}

// A measured run over the 24h ceiling can't be logged as-is, so stop must refuse
// early (before any network) and point at --duration, keeping the timer.
func TestStopRefusesOver24hWithoutDuration(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))
	mustStart(t, timer.State{TicketCode: "APP-2", ActivityName: "dev", CompanyID: 1, StartedAt: time.Now().Add(-25 * time.Hour)})

	_, _, err := runCmd(t, "stop", "--yes")
	if err == nil {
		t.Fatal("expected stop to refuse a >24h measured timer")
	}
	if !strings.Contains(err.Error(), "--duration") {
		t.Errorf("error = %q, want it to point at --duration", err.Error())
	}
	if _, loadErr := timer.Load(); loadErr != nil {
		t.Fatalf("timer should be kept after refusal, Load = %v", loadErr)
	}
}

// A corrupt timer file must be surfaced with the tl-cancel guidance up front,
// before any network work — not masked as "already running".
func TestStartOnCorruptTimerPointsToCancel(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TL_CONFIG_DIR", dir)
	mustSaveToken(t, time.Now().Add(time.Hour))
	if err := os.WriteFile(filepath.Join(dir, "timer.json"), []byte("{bad"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, _, err := runCmd(t, "start", "APP-1", "-a", "dev", "--base-url", "http://127.0.0.1:0", "--company", "1")
	if err == nil {
		t.Fatal("expected start to fail on a corrupt timer file")
	}
	if !strings.Contains(err.Error(), "tl cancel") {
		t.Errorf("error = %q, want it to suggest 'tl cancel'", err.Error())
	}
}

func TestCancelRemovesTimer(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustStart(t, timer.State{TicketCode: "APP-3", ActivityName: "dev", CompanyID: 1, StartedAt: time.Now().Add(-5 * time.Minute)})

	if _, _, err := runCmd(t, "cancel", "--yes"); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if _, err := timer.Load(); !errors.Is(err, timer.ErrNotRunning) {
		t.Fatalf("timer should be gone after cancel, Load = %v", err)
	}
}

func mustStart(t *testing.T, s timer.State) {
	t.Helper()
	if err := timer.Start(s); err != nil {
		t.Fatalf("seeding timer: %v", err)
	}
}

func mustSaveToken(t *testing.T, expiresAt time.Time) {
	t.Helper()
	if err := auth.Save(auth.Token{AccessToken: "test-token", ExpiresAt: expiresAt}); err != nil {
		t.Fatalf("saving token: %v", err)
	}
}
