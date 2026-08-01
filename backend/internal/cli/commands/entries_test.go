package commands

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/auth"
)

// entryServer scripts the endpoints `tl edit` and `tl delete` touch for entry
// #421: the GET/PUT/DELETE on the entry plus the timesheet lookup that resolves
// its ticket code. notFound flips the entry endpoints to 404.
type entryServer struct {
	notFound bool

	putBody      []byte
	putCalled    bool
	deleteCalled bool
}

func (s *entryServer) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/timesheet", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("start")
		writeJSON(w, http.StatusOK, map[string]any{
			"consultant_id": 1, "company_id": 1, "start": date, "end": date,
			"rows": []any{map[string]any{
				"ticket":   map[string]any{"id": 10, "company_id": 1, "code": "APP-123"},
				"activity": map[string]any{"id": 5, "company_id": 1, "name": "dev", "priority": 1},
				"entries": []any{map[string]any{
					"id": 421, "date": date, "duration_minutes": 45, "comment": "old",
				}},
				"per_day_minutes": map[string]int{date: 45},
				"total_minutes":   45,
			}},
			"totals": map[string]any{"per_day_minutes": map[string]int{date: 45}, "overall_minutes": 45},
		})
	})

	mux.HandleFunc("/api/entries/421", func(w http.ResponseWriter, r *http.Request) {
		if s.notFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "entry not found"})
			return
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, map[string]any{
				"id": 421, "consultant_id": 1, "company_id": 1, "contract_id": 1,
				"ticket_id": 10, "activity_id": 5, "date": "2026-07-13",
				"duration_minutes": 45, "comment": "old",
			})
		case http.MethodPut:
			s.putCalled = true
			s.putBody, _ = io.ReadAll(r.Body)
			writeJSON(w, http.StatusOK, map[string]any{
				"id": 421, "consultant_id": 1, "company_id": 1, "contract_id": 1,
				"ticket_id": 10, "activity_id": 5, "date": "2026-07-13",
				"duration_minutes": 120, "comment": "old",
			})
		case http.MethodDelete:
			s.deleteCalled = true
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}

// runAuthedCmd executes a subcommand against srv with a valid token in a temp
// config dir. stdin is the given string (empty = non-interactive EOF for prompts).
func runAuthedCmd(t *testing.T, srvURL, stdin string, args ...string) (string, error) {
	t.Helper()
	t.Setenv("TL_CONFIG_DIR", t.TempDir())
	if err := auth.Save(auth.Token{AccessToken: "test-token", ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("saving token: %v", err)
	}

	var out bytes.Buffer
	cmd := NewRoot("test")
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)
	cmd.SetIn(strings.NewReader(stdin))
	full := append([]string{args[0], "--base-url", srvURL, "--company", "1"}, args[1:]...)
	cmd.SetArgs(full)
	err := cmd.Execute()
	return out.String(), err
}
