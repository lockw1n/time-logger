package outbox

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestSyncLockContention: two Sync calls racing over the same outbox must result
// in each op being replayed exactly once — the lock guarantees a single replayer.
func TestSyncLockContention(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	var mu sync.Mutex
	posts := 0
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		posts++
		mu.Unlock()
		<-release // hold the winner inside the locked region so the loser races in
		writeJSON(w, http.StatusCreated, map[string]any{"id": posts, "date": "2026-08-01"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	enqueueOp(t, "APP-1", 30)
	enqueueOp(t, "APP-2", 45)

	client := newClient(srv.URL)
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Go(func() {
			_, err := Sync(context.Background(), client, Options{})
			results <- err
		})
	}

	// Let the winner make progress, then release both requests and finish.
	time.Sleep(100 * time.Millisecond)
	close(release)
	wg.Wait()
	close(results)

	locked := 0
	for err := range results {
		if errors.Is(err, ErrLocked) {
			locked++
		} else if err != nil {
			t.Fatalf("unexpected Sync error: %v", err)
		}
	}
	if locked != 1 {
		t.Fatalf("exactly one Sync should be locked out, got %d", locked)
	}

	mu.Lock()
	got := posts
	mu.Unlock()
	if got != 2 {
		t.Fatalf("server saw %d POSTs, want 2 (each op replayed once)", got)
	}
	if n, _ := Count(); n != 0 {
		t.Fatalf("outbox should be empty, has %d", n)
	}
}

// TestSyncBreaksStaleLock: a lock left by a dead process must be broken (with a
// notice) so a crashed sync can't wedge the outbox forever.
func TestSyncBreaksStaleLock(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TL_CONFIG_DIR", dir)

	// Pretend every recorded pid is dead.
	orig := processAlive
	processAlive = func(int) bool { return false }
	defer func() { processAlive = orig }()

	outboxDir := filepath.Join(dir, dirName)
	if err := os.MkdirAll(outboxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outboxDir, lockName), []byte("999999\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, map[string]any{"id": 1, "date": "2026-08-01"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	enqueueOp(t, "APP-1", 30)

	var notice bytes.Buffer
	rep, err := Sync(context.Background(), newClient(srv.URL), Options{Notice: &notice})
	if err != nil {
		t.Fatalf("Sync should break the stale lock and run, got %v", err)
	}
	if rep.Synced != 1 {
		t.Fatalf("Synced = %d, want 1", rep.Synced)
	}
	if !strings.Contains(notice.String(), "stale") {
		t.Errorf("notice = %q, want it to mention breaking the stale lock", notice.String())
	}
}

// TestSyncRespectsLiveLock: a lock held by a live process yields ErrLocked and
// leaves the queue untouched.
func TestSyncRespectsLiveLock(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TL_CONFIG_DIR", dir)

	outboxDir := filepath.Join(dir, dirName)
	if err := os.MkdirAll(outboxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	// Our own pid is unquestionably alive.
	if err := os.WriteFile(filepath.Join(outboxDir, lockName), []byte("1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Force the recorded pid to look alive regardless of platform quirks.
	orig := processAlive
	processAlive = func(int) bool { return true }
	defer func() { processAlive = orig }()

	enqueueOp(t, "APP-1", 30)

	rep, err := Sync(context.Background(), newClient("http://127.0.0.1:0"), Options{})
	if !errors.Is(err, ErrLocked) {
		t.Fatalf("Sync err = %v, want ErrLocked", err)
	}
	if rep.Synced != 0 {
		t.Fatalf("Synced = %d, want 0 while locked", rep.Synced)
	}
	if n, _ := Count(); n != 1 {
		t.Fatalf("op should remain queued while locked, Count = %d", n)
	}
}
