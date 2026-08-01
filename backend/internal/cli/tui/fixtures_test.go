package tui

import "github.com/lockw1n/time-logger/internal/cli/api"

func ptr(s string) *string { return &s }

// weekFixture is a full Monday–Sunday (2026-07-06 … 2026-07-12) timesheet with
// two rows: one with a deliberately long activity name (to exercise truncation)
// and a day carrying two entries (to exercise the multi-entry cell path). Rows
// are supplied out of display order so the sort is exercised too.
func weekFixture() api.TimesheetResponse {
	return api.TimesheetResponse{
		Start: "2026-07-06",
		End:   "2026-07-12",
		Rows: []api.TimesheetRow{
			{
				Ticket:   api.TimesheetTicket{Code: "APP-200"},
				Activity: api.TimesheetActivity{ID: 2, Name: "review", Priority: 2},
				Entries: []api.TimesheetEntry{
					{ID: 55, Date: "2026-07-08", DurationMinutes: 45, Comment: ptr("morning pass")},
					{ID: 60, Date: "2026-07-08", DurationMinutes: 60, Comment: ptr("second pass")},
				},
				PerDayMinutes: map[string]int{"2026-07-08": 105},
				TotalMinutes:  105,
			},
			{
				Ticket:   api.TimesheetTicket{Code: "APP-123"},
				Activity: api.TimesheetActivity{ID: 1, Name: "Spryker Feature Development", Priority: 1},
				Entries: []api.TimesheetEntry{
					{ID: 42, Date: "2026-07-06", DurationMinutes: 120},
					{ID: 43, Date: "2026-07-07", DurationMinutes: 90},
					{ID: 44, Date: "2026-07-09", DurationMinutes: 180},
				},
				PerDayMinutes: map[string]int{"2026-07-06": 120, "2026-07-07": 90, "2026-07-09": 180},
				TotalMinutes:  390,
			},
		},
		Totals: api.TimesheetTotals{
			PerDayMinutes: map[string]int{
				"2026-07-06": 120, "2026-07-07": 90, "2026-07-08": 105, "2026-07-09": 180,
			},
			OverallMinutes: 495,
		},
	}
}
