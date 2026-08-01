package timer

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func setupDir(t *testing.T) {
	t.Helper()
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
}

func TestStartLoadClearLifecycle(t *testing.T) {
	setupDir(t)

	if _, err := Load(); !errors.Is(err, ErrNotRunning) {
		t.Fatalf("Load before start = %v, want ErrNotRunning", err)
	}

	started := time.Date(2026, 7, 13, 9, 15, 0, 0, time.UTC)
	want := State{TicketCode: "APP-123", ActivityName: "dev", CompanyID: 1, StartedAt: started, Note: "hi"}
	if err := Start(want); err != nil {
		t.Fatalf("Start: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.TicketCode != want.TicketCode || got.ActivityName != want.ActivityName ||
		got.CompanyID != want.CompanyID || !got.StartedAt.Equal(want.StartedAt) || got.Note != want.Note {
		t.Fatalf("Load = %+v, want %+v", got, want)
	}

	// A second start while one runs must be refused.
	if err := Start(want); !errors.Is(err, ErrRunning) {
		t.Fatalf("second Start = %v, want ErrRunning", err)
	}

	if err := Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := Load(); !errors.Is(err, ErrNotRunning) {
		t.Fatalf("Load after clear = %v, want ErrNotRunning", err)
	}
	// Clearing an absent timer is a no-op.
	if err := Clear(); err != nil {
		t.Fatalf("Clear when absent: %v", err)
	}
}

// Two concurrent starts must resolve to exactly one winner via the atomic
// O_EXCL create — never two timers, never zero.
func TestStartConcurrentExactlyOneWins(t *testing.T) {
	setupDir(t)

	const n = 8
	var wg sync.WaitGroup
	var wins int64
	started := time.Now()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := Start(State{TicketCode: "APP-1", ActivityName: "dev", CompanyID: 1, StartedAt: started})
			switch {
			case err == nil:
				atomic.AddInt64(&wins, 1)
			case errors.Is(err, ErrRunning):
				// expected loser
			default:
				t.Errorf("goroutine %d: unexpected error %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	if wins != 1 {
		t.Fatalf("winners = %d, want exactly 1", wins)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	setupDir(t)

	path, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err = Load()
	var corrupt *CorruptError
	if !errors.As(err, &corrupt) {
		t.Fatalf("Load = %v, want *CorruptError", err)
	}
	if !strings.Contains(err.Error(), "tl cancel") {
		t.Errorf("corrupt error %q should suggest 'tl cancel'", err.Error())
	}
}

// A well-formed JSON object that is missing required fields is still corrupt: a
// truncated write could leave valid JSON without a ticket or start time.
func TestLoadMissingFieldsIsCorrupt(t *testing.T) {
	setupDir(t)

	path, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"activity_name":"dev"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err = Load()
	var corrupt *CorruptError
	if !errors.As(err, &corrupt) {
		t.Fatalf("Load = %v, want *CorruptError", err)
	}
}

func TestElapsedMinutes(t *testing.T) {
	started := time.Date(2026, 7, 13, 9, 0, 0, 0, time.UTC)
	s := State{StartedAt: started}

	cases := []struct {
		name string
		now  time.Time
		want int
	}{
		{"47 minutes", started.Add(47 * time.Minute), 47},
		{"zero floors to 1", started, 1},
		{"29s rounds to 0, floored to 1", started.Add(29 * time.Second), 1},
		{"90s rounds to 2", started.Add(90 * time.Second), 2},
		{"exactly two hours", started.Add(2 * time.Hour), 120},
	}
	for _, c := range cases {
		if got := ElapsedMinutes(s.Elapsed(c.now)); got != c.want {
			t.Errorf("%s: ElapsedMinutes = %d, want %d", c.name, got, c.want)
		}
	}
}
