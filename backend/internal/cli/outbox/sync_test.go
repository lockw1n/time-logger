package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// newClient builds a real API client pointed at srv with a static token, so the
// tests exercise the genuine create/merge path (entrylog.Submit) end to end.
func newClient(url string) api.Client {
	return api.New(url, func() (string, error) { return "test-token", nil })
}

// enqueueOp queues a create for the given ticket/minutes and returns nothing —
// filenames are time-ordered, so enqueue order is replay order.
func enqueueOp(t *testing.T, ticket string, minutes int) {
	t.Helper()
	if _, err := Enqueue(Op{
		Kind:      KindCreateEntry,
		CreatedAt: time.Now(),
		Payload: api.CreateEntryRequest{
			CompanyID:       1,
			TicketCode:      ticket,
			ActivityID:      5,
			Date:            "2026-08-01",
			DurationMinutes: minutes,
		},
		MergeOnConflict: true,
	}); err != nil {
		t.Fatalf("enqueue %s: %v", ticket, err)
	}
	// Guarantee a distinct nanosecond prefix so ordering is unambiguous even on
	// coarse clocks.
	time.Sleep(time.Millisecond)
}

// TestSyncReplaysOldestFirst verifies ops are replayed in enqueue order and each
// synced op's file is removed.
func TestSyncReplaysOldestFirst(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	var mu sync.Mutex
	var order []string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		var req api.CreateEntryRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		mu.Lock()
		order = append(order, req.TicketCode)
		mu.Unlock()
		writeJSON(w, http.StatusCreated, map[string]any{"id": len(order), "date": "2026-08-01"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	enqueueOp(t, "APP-1", 30)
	enqueueOp(t, "APP-2", 45)
	enqueueOp(t, "APP-3", 60)

	rep, err := Sync(context.Background(), newClient(srv.URL), Options{})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if rep.Synced != 3 || rep.Remaining != 0 || rep.Failed != 0 {
		t.Fatalf("report = %+v, want {Synced:3, Failed:0, Remaining:0}", rep)
	}
	want := []string{"APP-1", "APP-2", "APP-3"}
	if len(order) != 3 || order[0] != want[0] || order[1] != want[1] || order[2] != want[2] {
		t.Fatalf("replay order = %v, want %v", order, want)
	}
	if n, _ := Count(); n != 0 {
		t.Fatalf("outbox should be empty after sync, has %d", n)
	}
}

// TestSyncStopsWhenUnreachableMidReplay: the server accepts one op then becomes
// unreachable. The first file must be gone, the rest untouched, Remaining right.
func TestSyncStopsWhenUnreachableMidReplay(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	var mu sync.Mutex
	calls := 0
	closeAfterFirst := make(chan struct{}, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			var req api.CreateEntryRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			writeJSON(w, http.StatusCreated, map[string]any{"id": 1, "date": "2026-08-01"})
			select {
			case closeAfterFirst <- struct{}{}:
			default:
			}
			return
		}
		// Should not be reached once the server is closed, but be safe.
		writeJSON(w, http.StatusCreated, map[string]any{"id": n, "date": "2026-08-01"})
	})
	srv := httptest.NewServer(mux)

	enqueueOp(t, "APP-1", 30)
	enqueueOp(t, "APP-2", 45)
	enqueueOp(t, "APP-3", 60)

	client := newClient(srv.URL)
	// Close the server as soon as the first op is accepted so the second dial is
	// refused (unreachable) — deterministically stopping the run after one op.
	go func() {
		<-closeAfterFirst
		srv.Close()
	}()

	rep, err := Sync(context.Background(), client, Options{})
	if err != nil {
		t.Fatalf("Sync should stop cleanly on unreachable, got err %v", err)
	}
	if rep.Synced != 1 {
		t.Fatalf("Synced = %d, want 1", rep.Synced)
	}
	if rep.Remaining != 2 {
		t.Fatalf("Remaining = %d, want 2", rep.Remaining)
	}
	if n, _ := Count(); n != 2 {
		t.Fatalf("outbox should retain 2 ops, has %d", n)
	}
}

// TestSyncMergesOnConflictWithoutPrompt: a queued op for an already-existing
// entry must merge (PUT with summed minutes) with no prompt, even with stdin
// closed — replay is non-interactive by construction.
func TestSyncMergesOnConflictWithoutPrompt(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	var putBody []byte
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "conflict"})
	})
	mux.HandleFunc("/api/timesheet", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("start")
		writeJSON(w, http.StatusOK, map[string]any{
			"consultant_id": 1, "company_id": 1, "start": date, "end": date,
			"rows": []any{map[string]any{
				"ticket":   map[string]any{"id": 10, "company_id": 1, "code": "APP-1"},
				"activity": map[string]any{"id": 5, "company_id": 1, "name": "dev", "priority": 1},
				"entries": []any{map[string]any{
					"id": 42, "date": date, "duration_minutes": 60, "comment": "earlier",
				}},
				"per_day_minutes": map[string]int{date: 60},
				"total_minutes":   60,
			}},
			"totals": map[string]any{"per_day_minutes": map[string]int{date: 60}, "overall_minutes": 60},
		})
	})
	mux.HandleFunc("/api/entries/42", func(w http.ResponseWriter, r *http.Request) {
		putBody, _ = io.ReadAll(r.Body)
		writeJSON(w, http.StatusOK, map[string]any{
			"id": 42, "date": "2026-08-01", "duration_minutes": 90, "comment": "earlier",
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	enqueueOp(t, "APP-1", 30) // 60 existing + 30 → 90

	rep, err := Sync(context.Background(), newClient(srv.URL), Options{})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if rep.Synced != 1 || rep.Remaining != 0 {
		t.Fatalf("report = %+v, want Synced 1 / Remaining 0", rep)
	}
	if putBody == nil {
		t.Fatal("expected a PUT to merge the entry, none sent")
	}
	var put struct {
		DurationMinutes *int `json:"duration_minutes"`
	}
	_ = json.Unmarshal(putBody, &put)
	if put.DurationMinutes == nil || *put.DurationMinutes != 90 {
		t.Fatalf("PUT duration = %v, want 90 (60 + 30)", put.DurationMinutes)
	}
}

// TestSyncQuarantinesPoisonOp: a 400 (validation) op can never succeed, so it is
// moved to failed/ with last_error and the run continues with the next op.
func TestSyncQuarantinesPoisonOp(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		var req api.CreateEntryRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.TicketCode == "BAD" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ticket"})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"id": 1, "date": "2026-08-01"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	enqueueOp(t, "BAD", 30)
	enqueueOp(t, "APP-2", 45)

	rep, err := Sync(context.Background(), newClient(srv.URL), Options{})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if rep.Synced != 1 || rep.Failed != 1 || rep.Remaining != 0 {
		t.Fatalf("report = %+v, want Synced 1 / Failed 1 / Remaining 0", rep)
	}
	failed, err := Failed()
	if err != nil {
		t.Fatalf("Failed(): %v", err)
	}
	if len(failed) != 1 {
		t.Fatalf("failed/ has %d ops, want 1", len(failed))
	}
	if failed[0].Op.Payload.TicketCode != "BAD" {
		t.Errorf("quarantined op = %q, want BAD", failed[0].Op.Payload.TicketCode)
	}
	if failed[0].Op.LastError == "" {
		t.Error("quarantined op should record a last_error")
	}
}

// TestSyncStopsPromptlyOnCancel: cancelling mid-run makes Sync return quickly
// with the in-flight (unconfirmed) op still queued — nothing is lost.
func TestSyncStopsPromptlyOnCancel(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	started := make(chan struct{})
	serverDone := make(chan struct{})
	var once sync.Once
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(started) })
		// Hang until the client cancels the request (the case under test) or the
		// test finishes, so srv.Close() never blocks on this handler.
		select {
		case <-r.Context().Done():
		case <-serverDone:
		}
	})
	srv := httptest.NewServer(mux)
	// Deferred LIFO: release the handler first, then close the server, so
	// srv.Close() never waits on a still-blocked request.
	defer srv.Close()
	defer close(serverDone)

	enqueueOp(t, "APP-1", 30)
	enqueueOp(t, "APP-2", 45)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan Report, 1)
	go func() {
		rep, _ := Sync(ctx, newClient(srv.URL), Options{})
		done <- rep
	}()

	<-started
	cancel()

	select {
	case rep := <-done:
		if rep.Remaining != 2 {
			t.Fatalf("Remaining = %d, want 2 (nothing lost on cancel)", rep.Remaining)
		}
		if rep.Synced != 0 {
			t.Fatalf("Synced = %d, want 0", rep.Synced)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Sync did not stop promptly after cancel")
	}
}

// TestSyncQuarantinesCorruptOpRaw: an unparseable op file must not be replayed
// as a zero payload; it is moved to failed/ byte-for-byte so its content survives.
func TestSyncQuarantinesCorruptOpRaw(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TL_CONFIG_DIR", dir)

	var posted bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		posted = true
		writeJSON(w, http.StatusCreated, map[string]any{"id": 1, "date": "2026-08-01"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// A valid op followed by a corrupt one.
	enqueueOp(t, "APP-1", 30)
	outboxDir := filepath.Join(dir, dirName)
	corrupt := []byte("{ this is not valid json")
	corruptPath := filepath.Join(outboxDir, fmt.Sprintf("%d-dead.json", time.Now().UnixNano()))
	if err := os.WriteFile(corruptPath, corrupt, 0o600); err != nil {
		t.Fatal(err)
	}

	rep, err := Sync(context.Background(), newClient(srv.URL), Options{})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if !posted {
		t.Error("the valid op should still have been replayed")
	}
	if rep.Synced != 1 || rep.Failed != 1 {
		t.Fatalf("report = %+v, want Synced 1 / Failed 1", rep)
	}

	// The corrupt file must be preserved verbatim in failed/, not a zeroed re-encode.
	failedPath := filepath.Join(dir, dirName, failedName, filepath.Base(corruptPath))
	got, err := os.ReadFile(failedPath)
	if err != nil {
		t.Fatalf("corrupt op should be quarantined at %s: %v", failedPath, err)
	}
	if string(got) != string(corrupt) {
		t.Errorf("quarantined bytes = %q, want the original %q", got, corrupt)
	}
}

// writeJSON is the test JSON helper mirroring the one in the commands package.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
