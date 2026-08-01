package parse

import (
	"testing"
	"time"
)

func TestMonthRange(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		start, end string
		wantErr    bool
	}{
		{name: "31-day month", in: "2026-01", start: "2026-01-01", end: "2026-01-31"},
		{name: "30-day month", in: "2026-06", start: "2026-06-01", end: "2026-06-30"},
		{name: "february non-leap", in: "2025-02", start: "2025-02-01", end: "2025-02-28"},
		{name: "february leap", in: "2024-02", start: "2024-02-01", end: "2024-02-29"},
		{name: "december", in: "2026-12", start: "2026-12-01", end: "2026-12-31"},
		{name: "whitespace trimmed", in: "  2026-06  ", start: "2026-06-01", end: "2026-06-30"},
		{name: "bad format full date", in: "2026-06-01", wantErr: true},
		{name: "bad month", in: "2026-13", wantErr: true},
		{name: "empty", in: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, end, err := MonthRange(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("MonthRange(%q) = (%q, %q, nil), want error", tc.in, start, end)
				}
				return
			}
			if err != nil {
				t.Fatalf("MonthRange(%q) error: %v", tc.in, err)
			}
			if start != tc.start || end != tc.end {
				t.Errorf("MonthRange(%q) = (%q, %q), want (%q, %q)", tc.in, start, end, tc.start, tc.end)
			}
		})
	}
}

func TestPreviousMonth(t *testing.T) {
	tests := []struct {
		name string
		now  time.Time
		want string
	}{
		{name: "mid-year", now: date(2026, 7, 31), want: "2026-06"},
		{name: "january rolls to previous december", now: date(2026, 1, 15), want: "2025-12"},
		{name: "march after leap february", now: date(2024, 3, 1), want: "2024-02"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := PreviousMonth(tc.now); got != tc.want {
				t.Errorf("PreviousMonth(%v) = %q, want %q", tc.now, got, tc.want)
			}
		})
	}
}

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 10, 0, 0, 0, time.UTC)
}
