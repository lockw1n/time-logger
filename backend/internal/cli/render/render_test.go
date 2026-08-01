package render

import (
	"testing"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

func ptr(s string) *string { return &s }

// multiRowDay is a single-day (2026-07-30) timesheet with two rows, one of them
// carrying a comment. Rows are supplied out of order to exercise sorting.
func multiRowDay() api.TimesheetResponse {
	return api.TimesheetResponse{
		Start: "2026-07-30",
		End:   "2026-07-30",
		Rows: []api.TimesheetRow{
			{
				Ticket:   api.TimesheetTicket{Code: "APP-200"},
				Activity: api.TimesheetActivity{Name: "review", Priority: 2},
				Entries: []api.TimesheetEntry{
					{ID: 55, Date: "2026-07-30", DurationMinutes: 45},
				},
				PerDayMinutes: map[string]int{"2026-07-30": 45},
				TotalMinutes:  45,
			},
			{
				Ticket:   api.TimesheetTicket{Code: "APP-123"},
				Activity: api.TimesheetActivity{Name: "dev", Priority: 1},
				Entries: []api.TimesheetEntry{
					{ID: 42, Date: "2026-07-30", DurationMinutes: 90, Comment: ptr("code review")},
				},
				PerDayMinutes: map[string]int{"2026-07-30": 90},
				TotalMinutes:  90,
			},
		},
		Totals: api.TimesheetTotals{
			PerDayMinutes:  map[string]int{"2026-07-30": 135},
			OverallMinutes: 135,
		},
	}
}

func TestTodayMultiRow(t *testing.T) {
	want := "" +
		"APP-123  dev     1h30m  code review  #42\n" +
		"APP-200  review  45m                 #55\n" +
		"total: 2h15m\n"

	if got := Today(multiRowDay()); got != want {
		t.Errorf("Today() mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTodayEmpty(t *testing.T) {
	ts := api.TimesheetResponse{Start: "2026-07-30", End: "2026-07-30"}
	want := "no entries for 2026-07-30\n"
	if got := Today(ts); got != want {
		t.Errorf("Today(empty) = %q, want %q", got, want)
	}
}

func TestWeekMultiRow(t *testing.T) {
	ts := api.TimesheetResponse{
		Start: "2026-07-27", // Monday
		End:   "2026-08-02",
		Rows: []api.TimesheetRow{
			{
				Ticket:   api.TimesheetTicket{Code: "APP-123"},
				Activity: api.TimesheetActivity{Name: "dev", Priority: 1},
				PerDayMinutes: map[string]int{
					"2026-07-27": 120, "2026-07-28": 90, "2026-07-30": 180,
				},
				TotalMinutes: 390,
			},
			{
				Ticket:   api.TimesheetTicket{Code: "APP-200"},
				Activity: api.TimesheetActivity{Name: "review", Priority: 2},
				PerDayMinutes: map[string]int{
					"2026-07-29": 45, "2026-07-31": 60,
				},
				TotalMinutes: 105,
			},
		},
		Totals: api.TimesheetTotals{
			PerDayMinutes: map[string]int{
				"2026-07-27": 120, "2026-07-28": 90, "2026-07-29": 45,
				"2026-07-30": 180, "2026-07-31": 60,
			},
			OverallMinutes: 495,
		},
	}

	want := "" +
		"ticket   activity  Mon  Tue    Wed  Thu  Fri  Sat  Sun  total\n" +
		"APP-123  dev       2h   1h30m  —    3h   —    —    —    6h30m\n" +
		"APP-200  review    —    —      45m  —    1h   —    —    1h45m\n" +
		"total              2h   1h30m  45m  3h   1h   —    —    8h15m\n"

	if got := Week(ts); got != want {
		t.Errorf("Week() mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestWeekEmpty(t *testing.T) {
	ts := api.TimesheetResponse{Start: "2026-07-27", End: "2026-08-02"}
	want := "no entries for the week of 2026-07-27\n"
	if got := Week(ts); got != want {
		t.Errorf("Week(empty) = %q, want %q", got, want)
	}
}
