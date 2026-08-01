package tui

import "testing"

func TestBuildGridSortsRows(t *testing.T) {
	g := buildGrid(weekFixture())

	if len(g.rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(g.rows))
	}
	// Sorted by ticket code: APP-123 before APP-200 (fixture supplies them reversed).
	if g.rows[0].ticketCode != "APP-123" || g.rows[1].ticketCode != "APP-200" {
		t.Errorf("rows out of order: %q, %q", g.rows[0].ticketCode, g.rows[1].ticketCode)
	}
}

func TestBuildGridDaysAndTotals(t *testing.T) {
	g := buildGrid(weekFixture())

	wantDays := [7]string{
		"2026-07-06", "2026-07-07", "2026-07-08",
		"2026-07-09", "2026-07-10", "2026-07-11", "2026-07-12",
	}
	if g.days != wantDays {
		t.Errorf("days = %v, want %v", g.days, wantDays)
	}

	// APP-123: Mon 120, Tue 90, Thu 180.
	app123 := g.rows[0]
	if app123.perDay[0] != 120 || app123.perDay[1] != 90 || app123.perDay[3] != 180 {
		t.Errorf("APP-123 perDay = %v", app123.perDay)
	}
	if app123.perDay[2] != 0 {
		t.Errorf("APP-123 Wed should be empty, got %d", app123.perDay[2])
	}

	if g.overall != 495 {
		t.Errorf("overall = %d, want 495", g.overall)
	}
	if g.totals[2] != 105 { // Wednesday
		t.Errorf("totals[Wed] = %d, want 105", g.totals[2])
	}
}

// TestEntriesAt is the multi-entry-cell resolution the enter/delete dispatch
// depends on: zero, one, and two-or-more entries in a cell.
func TestEntriesAt(t *testing.T) {
	g := buildGrid(weekFixture())
	// rows[0]=APP-123, rows[1]=APP-200; days index Mon=0..Sun=6, Wed=2, Thu=3.

	tests := []struct {
		name     string
		row, day int
		want     int
	}{
		{"empty cell", 0, 2, 0},       // APP-123 Wednesday
		{"single entry", 0, 0, 1},     // APP-123 Monday (#42)
		{"multi entry", 1, 2, 2},      // APP-200 Wednesday (#55,#60)
		{"out of range row", 9, 0, 0}, // guarded
		{"out of range day", 0, 9, 0}, // guarded
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(g.entriesAt(tt.row, tt.day)); got != tt.want {
				t.Errorf("entriesAt(%d,%d) = %d entries, want %d", tt.row, tt.day, got, tt.want)
			}
		})
	}
}

func TestWeekDaysFallback(t *testing.T) {
	// An unparseable start degrades to the same string in every slot rather than
	// panicking, so the renderer can still draw an (empty) frame.
	days := weekDays("not-a-date")
	for i, d := range days {
		if d != "not-a-date" {
			t.Errorf("days[%d] = %q, want fallback", i, d)
		}
	}
}
