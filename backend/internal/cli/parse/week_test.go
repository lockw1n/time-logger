package parse

import (
	"strings"
	"testing"
)

func TestWeekBounds(t *testing.T) {
	cases := []struct {
		in, mon, sun string
	}{
		{"2026-07-30", "2026-07-27", "2026-08-02"}, // Thursday
		{"2026-07-27", "2026-07-27", "2026-08-02"}, // Monday itself
		{"2026-08-02", "2026-07-27", "2026-08-02"}, // Sunday itself
	}
	for _, tc := range cases {
		mon, sun, err := WeekBounds(tc.in)
		if err != nil {
			t.Fatalf("WeekBounds(%q) error: %v", tc.in, err)
		}
		if mon != tc.mon || sun != tc.sun {
			t.Errorf("WeekBounds(%q) = %q..%q, want %q..%q", tc.in, mon, sun, tc.mon, tc.sun)
		}
	}

	if _, _, err := WeekBounds("nope"); err == nil {
		t.Error("WeekBounds(invalid) want error")
	}
}

func TestWeekDays(t *testing.T) {
	days, err := WeekDays("2026-07-30")
	if err != nil {
		t.Fatalf("WeekDays error: %v", err)
	}
	want := "2026-07-27 2026-07-28 2026-07-29 2026-07-30 2026-07-31 2026-08-01 2026-08-02"
	if got := strings.Join(days, " "); got != want {
		t.Errorf("WeekDays = %q, want %q", got, want)
	}
}
