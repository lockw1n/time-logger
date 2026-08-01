package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/outbox"
)

// TestAddQueuesWhenUnreachable covers the headline acceptance criterion: with the
// backend down, `tl add` resolves the activity from the offline cache, queues the
// create, prints the queued message, and exits 3 — losing nothing.
func TestAddQueuesWhenUnreachable(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))
	// Warm the activities cache as a prior online command would have.
	saveActivitiesCache(1, []api.Activity{{ID: 5, CompanyID: 1, Name: "dev"}})

	// A server that is immediately closed → connection refused.
	ts := httptest.NewServer(activitiesOnlyMux())
	url := ts.URL
	ts.Close()

	out, _, err := runCmd(t, "add", "APP-123", "45m", "--activity", "dev", "--yes", "--base-url", url, "--company", "1")
	if got := exitCode(t, err); got != 3 {
		t.Fatalf("add exit = %d, want 3 (queued)", got)
	}
	if !strings.Contains(out, "queued") || !strings.Contains(out, "APP-123") {
		t.Errorf("stdout = %q, want the queued message naming APP-123", out)
	}

	pending, perr := outbox.Pending()
	if perr != nil {
		t.Fatalf("reading outbox: %v", perr)
	}
	if len(pending) != 1 {
		t.Fatalf("outbox has %d ops, want 1", len(pending))
	}
	op := pending[0].Op
	if op.Payload.TicketCode != "APP-123" || op.Payload.ActivityID != 5 || op.Payload.DurationMinutes != 45 {
		t.Errorf("queued payload = %+v, want APP-123/activity 5/45m", op.Payload)
	}
	if !op.MergeOnConflict {
		t.Error("queued create should be merge-on-conflict")
	}
}

// TestAddColdCacheOfflineGivesGuidance: offline add with no cached activities
// can't build a payload to queue, so it must fail with actionable guidance
// (pointing at `tl activities`) rather than a bare dial error, and queue nothing.
func TestAddColdCacheOfflineGivesGuidance(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))
	// No saveActivitiesCache call → cold cache.

	ts := httptest.NewServer(activitiesOnlyMux())
	url := ts.URL
	ts.Close()

	_, _, err := runCmd(t, "add", "APP-1", "45m", "--activity", "dev", "--yes", "--base-url", url, "--company", "1")
	if err == nil {
		t.Fatal("expected an error when offline with a cold activities cache")
	}
	if !strings.Contains(err.Error(), "tl activities") {
		t.Errorf("error = %q, want it to point at 'tl activities'", err.Error())
	}
	if n, _ := outbox.Count(); n != 0 {
		t.Fatalf("nothing should be queued when the activity can't resolve, got %d", n)
	}
}

// TestReportSyncPartialRunShowsBoth: a run that syncs some ops and leaves others
// must report both facts, not just the leftover (regression for the exclusive
// switch that previously hid the synced count).
func TestReportSyncPartialRunShowsBoth(t *testing.T) {
	var buf strings.Builder
	r := &root{errOut: &buf}
	r.reportSync(outbox.Report{Synced: 2, Failed: 1, Remaining: 1}, nil, "")

	out := buf.String()
	for _, want := range []string{"synced 2 queued entries", "1 queued entry could not be synced", "sync incomplete — 1 op still queued"} {
		if !strings.Contains(out, want) {
			t.Errorf("report %q missing %q", out, want)
		}
	}
}

// TestAutoSyncOnCommandStartup: a queued op is replayed in the background when the
// next API command runs, and the one-line report is printed to stderr.
func TestAutoSyncOnCommandStartup(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))

	if _, err := outbox.Enqueue(outbox.Op{
		Kind:      outbox.KindCreateEntry,
		CreatedAt: time.Now(),
		Payload: api.CreateEntryRequest{
			CompanyID: 1, TicketCode: "APP-9", ActivityID: 5,
			Date: "2026-08-01", DurationMinutes: 30,
		},
		MergeOnConflict: true,
	}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	var posted bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		posted = true
		writeJSON(w, http.StatusCreated, map[string]any{"id": 1, "date": "2026-08-01"})
	})
	mux.HandleFunc("/api/timesheet", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("start")
		writeJSON(w, http.StatusOK, map[string]any{
			"consultant_id": 1, "company_id": 1, "start": date, "end": date,
			"rows":   []any{},
			"totals": map[string]any{"per_day_minutes": map[string]int{}, "overall_minutes": 0},
		})
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	_, errBuf, err := runCmd(t, "today", "--base-url", ts.URL, "--company", "1")
	if err != nil {
		t.Fatalf("today: %v", err)
	}
	if !posted {
		t.Fatal("background sync did not replay the queued op")
	}
	if !strings.Contains(errBuf, "synced 1 queued entry") {
		t.Errorf("stderr = %q, want it to report the sync", errBuf)
	}
	if n, _ := outbox.Count(); n != 0 {
		t.Fatalf("outbox should be empty after auto-sync, has %d", n)
	}
}

// TestSyncCommandListReplayDiscard exercises the manual `tl sync` surface.
func TestSyncCommandListReplayDiscard(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	mustSaveToken(t, time.Now().Add(time.Hour))

	// One good op and one that the server will reject as invalid (→ failed/).
	enqueueForSync(t, "APP-OK", 5)
	enqueueForSync(t, "BAD", 5)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		var req api.CreateEntryRequest
		_ = decodeJSON(r, &req)
		if req.TicketCode == "BAD" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ticket"})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"id": 1, "date": "2026-08-01"})
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// --list before replaying shows both pending, no failed yet.
	listOut, _, err := runCmd(t, "sync", "--list")
	if err != nil {
		t.Fatalf("sync --list: %v", err)
	}
	if !strings.Contains(listOut, "pending (2)") {
		t.Errorf("list = %q, want 2 pending", listOut)
	}

	// Replay: one syncs, one is quarantined → exit 1 (something failed).
	syncOut, _, err := runCmd(t, "sync", "--base-url", ts.URL, "--company", "1")
	if got := exitCode(t, err); got != 1 {
		t.Fatalf("sync exit = %d, want 1 (a failed op remains)", got)
	}
	if !strings.Contains(syncOut, "synced 1") || !strings.Contains(syncOut, "failed 1") {
		t.Errorf("sync out = %q, want synced 1 / failed 1", syncOut)
	}

	// The failed op is now inspectable; discard it after confirmation.
	failed, _ := outbox.Failed()
	if len(failed) != 1 {
		t.Fatalf("failed has %d ops, want 1", len(failed))
	}
	name := failed[0].Filename

	if _, _, err := runCmdStdin(t, "y\n", "sync", "--discard", name); err != nil {
		t.Fatalf("sync --discard: %v", err)
	}
	if failed, _ := outbox.Failed(); len(failed) != 0 {
		t.Fatalf("failed op should be gone after discard, %d remain", len(failed))
	}
}

// enqueueForSync queues a create op for the given ticket in the current config dir.
func enqueueForSync(t *testing.T, ticket string, minutes int) {
	t.Helper()
	if _, err := outbox.Enqueue(outbox.Op{
		Kind:      outbox.KindCreateEntry,
		CreatedAt: time.Now(),
		Payload: api.CreateEntryRequest{
			CompanyID: 1, TicketCode: ticket, ActivityID: 5,
			Date: "2026-08-01", DurationMinutes: minutes,
		},
		MergeOnConflict: true,
	}); err != nil {
		t.Fatalf("enqueue %s: %v", ticket, err)
	}
	time.Sleep(time.Millisecond) // distinct nanosecond prefix → stable ordering
}

// runCmdStdin runs tl with args and a canned stdin, returning stdout, stderr, err.
func runCmdStdin(t *testing.T, stdin string, args ...string) (string, string, error) {
	t.Helper()
	var out, errBuf strings.Builder
	cmd := NewRoot("test")
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), errBuf.String(), err
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
