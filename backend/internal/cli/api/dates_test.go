package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCanonicalDate(t *testing.T) {
	cases := map[string]string{
		"30.07.2026": "2026-07-30", // response format -> canonical
		"2026-07-30": "2026-07-30", // already canonical, untouched
		"":           "",           // empty untouched
		"garbage":    "garbage",    // unparseable untouched
	}
	for in, want := range cases {
		if got := canonicalDate(in); got != want {
			t.Errorf("canonicalDate(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestTimesheetNormalizesDates verifies the client rewrites the DD.MM.YYYY
// scalar date fields the backend really returns into canonical YYYY-MM-DD.
func TestTimesheetNormalizesDates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"consultant_id":1,"company_id":1,"start":"27.07.2026","end":"02.08.2026",
			"rows":[{"ticket":{"code":"APP-1"},"activity":{"id":5,"name":"dev"},
				"entries":[{"id":9,"date":"30.07.2026","duration_minutes":60}],
				"per_day_minutes":{"2026-07-30":60},"total_minutes":60}],
			"totals":{"per_day_minutes":{"2026-07-30":60},"overall_minutes":60}
		}`))
	}))
	defer server.Close()

	client := New(server.URL, func() (string, error) { return "t", nil })
	ts, err := client.Timesheet(context.Background(), 1, "2026-07-27", "2026-08-02")
	if err != nil {
		t.Fatalf("Timesheet error: %v", err)
	}

	if ts.Start != "2026-07-27" || ts.End != "2026-08-02" {
		t.Errorf("start/end = %q/%q, want canonical", ts.Start, ts.End)
	}
	if got := ts.Rows[0].Entries[0].Date; got != "2026-07-30" {
		t.Errorf("entry date = %q, want 2026-07-30", got)
	}
	if _, ok := ts.Rows[0].PerDayMinutes["2026-07-30"]; !ok {
		t.Error("per_day_minutes key should remain canonical 2026-07-30")
	}
}
