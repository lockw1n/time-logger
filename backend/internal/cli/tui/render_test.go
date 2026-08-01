package tui

import "testing"

// TestRenderGridWide renders the fixture week at a width that fits every column
// (the full table is 96 cols), so the long activity name is shown in full. The
// cursor sits on the first row's Monday cell, which must appear framed in
// [brackets]; the Wednesday column (today) carries the leading '*' marker.
func TestRenderGridWide(t *testing.T) {
	g := buildGrid(weekFixture())
	todayCol := todayColumn(g, "2026-07-08") // Wednesday
	got := renderGrid(g, cursor{row: 0, day: 0}, todayCol, 120)

	want := "" +
		"ticket   activity                      Mon    Tue     *Wed      Thu    Fri    Sat    Sun   total\n" +
		"APP-123  Spryker Feature Development  [2h ]   1h30m    —        3h     —      —      —     6h30m\n" +
		"APP-200  review                        —      —        1h45m    —      —      —      —     1h45m\n" +
		"total                                  2h     1h30m    1h45m    3h     —      —      —     8h15m\n"

	if got != want {
		t.Errorf("renderGrid(wide) mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestRenderGridNarrow renders the same week at a width that forces the activity
// column to truncate with an ellipsis, proving the grid adapts to the terminal
// rather than overflowing. The cursor is on the second row's Thursday (empty)
// cell, which must still show framed brackets.
func TestRenderGridNarrow(t *testing.T) {
	g := buildGrid(weekFixture())
	todayCol := todayColumn(g, "2026-07-08")
	got := renderGrid(g, cursor{row: 1, day: 3}, todayCol, 80)

	want := "" +
		"ticket   activity      Mon    Tue     *Wed      Thu    Fri    Sat    Sun   total\n" +
		"APP-123  Spryker Fe…   2h     1h30m    —        3h     —      —      —     6h30m\n" +
		"APP-200  review        —      —        1h45m   [—  ]   —      —      —     1h45m\n" +
		"total                  2h     1h30m    1h45m    3h     —      —      —     8h15m\n"

	if got != want {
		t.Errorf("renderGrid(narrow) mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestRenderGridEmpty covers the no-entries week.
func TestRenderGridEmpty(t *testing.T) {
	g := buildGrid(weekFixture())
	g.rows = nil
	want := "no entries for the week of 2026-07-06\n"
	if got := renderGrid(g, cursor{}, -1, 120); got != want {
		t.Errorf("renderGrid(empty) = %q, want %q", got, want)
	}
}
