package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/auth"
)

// addTestServer scripts the endpoints `tl add` touches. matchDuplicate controls
// whether the timesheet lookup after the 409 reports an existing entry (merge)
// or nothing (contract error). Any PUT body is captured for assertions.
type addTestServer struct {
	matchDuplicate  bool
	existingMinutes int // minutes on the existing entry (defaults to 90)
	putBody         []byte
	putCalled       bool
}

func (s *addTestServer) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/activities", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, []map[string]any{
			{"id": 5, "company_id": 1, "name": "dev", "billable": true, "priority": 1},
		})
	})

	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, _ *http.Request) {
		// Always conflict so the merge-detection flow runs.
		writeJSON(w, http.StatusConflict, map[string]string{"error": "conflict"})
	})

	mux.HandleFunc("/api/timesheet", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("start")
		resp := map[string]any{
			"consultant_id": 1, "company_id": 1, "start": date, "end": date,
			"rows":   []any{},
			"totals": map[string]any{"per_day_minutes": map[string]int{}, "overall_minutes": 0},
		}
		if s.matchDuplicate {
			existing := s.existingMinutes
			if existing == 0 {
				existing = 90
			}
			resp["rows"] = []any{map[string]any{
				"ticket":   map[string]any{"id": 10, "company_id": 1, "code": "APP-123"},
				"activity": map[string]any{"id": 5, "company_id": 1, "name": "dev", "priority": 1},
				"entries": []any{map[string]any{
					"id": 42, "date": date, "duration_minutes": existing, "comment": "existing",
				}},
				"per_day_minutes": map[string]int{date: existing},
				"total_minutes":   existing,
			}}
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("/api/entries/42", func(w http.ResponseWriter, r *http.Request) {
		s.putCalled = true
		s.putBody, _ = io.ReadAll(r.Body)
		writeJSON(w, http.StatusOK, map[string]any{
			"id": 42, "company_id": 1, "activity_id": 5, "date": "2026-07-30",
			"duration_minutes": 135, "comment": "existing; new note",
		})
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// runAdd executes the add command against srv with a valid token in a temp
// config dir, returning stdout and the command error. Stdin is empty, so an
// interactive prompt reached without --yes sees immediate EOF.
func runAdd(t *testing.T, srvURL string, args ...string) (string, error) {
	t.Helper()
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	if err := auth.Save(auth.Token{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("saving token: %v", err)
	}

	var out bytes.Buffer
	cmd := NewRoot("test")
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)
	cmd.SetIn(strings.NewReader("")) // non-interactive: prompts see EOF
	cmd.SetArgs(append([]string{"add", "--base-url", srvURL, "--company", "1"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func TestAddMergeSumsMinutesAndComment(t *testing.T) {
	srv := &addTestServer{matchDuplicate: true}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	out, err := runAdd(t, ts.URL, "APP-123", "45m", "new note", "--activity", "dev", "--yes")
	if err != nil {
		t.Fatalf("add returned error: %v", err)
	}
	if !srv.putCalled {
		t.Fatal("expected a PUT to merge the entry, none was sent")
	}

	var put struct {
		DurationMinutes *int    `json:"duration_minutes"`
		Comment         *string `json:"comment"`
	}
	if err := json.Unmarshal(srv.putBody, &put); err != nil {
		t.Fatalf("decoding PUT body: %v", err)
	}
	if put.DurationMinutes == nil || *put.DurationMinutes != 135 {
		t.Errorf("PUT duration = %v, want 135 (90 + 45)", put.DurationMinutes)
	}
	if put.Comment == nil || *put.Comment != "existing; new note" {
		t.Errorf("PUT comment = %v, want %q", put.Comment, "existing; new note")
	}
	if !strings.Contains(out, "#42") {
		t.Errorf("output %q should reference merged entry #42", out)
	}
}

func TestAddNoContractExitsWithoutPut(t *testing.T) {
	srv := &addTestServer{matchDuplicate: false}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAdd(t, ts.URL, "APP-123", "45m", "--activity", "dev", "--yes")
	if err == nil {
		t.Fatal("expected a no-active-contract error, got nil")
	}
	if !strings.Contains(err.Error(), "no active contract") {
		t.Errorf("error = %q, want it to mention 'no active contract'", err.Error())
	}
	if srv.putCalled {
		t.Error("no PUT should be sent when the 409 is a contract problem")
	}
}

// A merge prompt reached over a non-interactive stdin (EOF) without --yes must
// abort rather than silently accepting the [Y/n] default.
func TestAddMergeAbortsOnNonInteractiveStdin(t *testing.T) {
	srv := &addTestServer{matchDuplicate: true}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAdd(t, ts.URL, "APP-123", "45m", "--activity", "dev") // no --yes
	if err == nil {
		t.Fatal("expected the merge to abort on EOF, got nil error")
	}
	if srv.putCalled {
		t.Error("no PUT should be sent when the merge prompt is not confirmed")
	}
}

// Merging past the 24h/entry ceiling must fail client-side, before any PUT.
func TestAddMergeRejectsOverDayLimit(t *testing.T) {
	srv := &addTestServer{matchDuplicate: true, existingMinutes: 23 * 60} // 23h already
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAdd(t, ts.URL, "APP-123", "90m", "--activity", "dev", "--yes")
	if err == nil {
		t.Fatal("expected a 24h-limit error, got nil")
	}
	if !strings.Contains(err.Error(), "24h") {
		t.Errorf("error = %q, want it to mention the 24h limit", err.Error())
	}
	if srv.putCalled {
		t.Error("no PUT should be sent when the merged total exceeds 24h")
	}
}
