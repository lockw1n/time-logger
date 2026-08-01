package parse

import (
	"errors"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	// Thursday, 2026-07-30, in UTC.
	now := time.Date(2026, 7, 30, 14, 0, 0, 0, time.UTC)

	cases := []struct {
		in      string
		want    string
		wantErr error
	}{
		{"", "2026-07-30", nil},
		{"today", "2026-07-30", nil},
		{"TODAY", "2026-07-30", nil},
		{"yesterday", "2026-07-29", nil},
		{"2026-01-15", "2026-01-15", nil},
		{"thu", "2026-07-30", nil}, // today is Thursday -> today
		{"wed", "2026-07-29", nil}, // most recent Wednesday
		{"fri", "2026-07-24", nil}, // last Friday (not the coming one)
		{"mon", "2026-07-27", nil},
		{"sun", "2026-07-26", nil},
		{"2026-08-01", "2026-08-01", ErrFutureDate},
		{"tomorrow", "", nil},   // not a recognized keyword
		{"2026-13-40", "", nil}, // impossible date
		{"07/30/2026", "", nil}, // wrong format
	}

	for _, tc := range cases {
		got, err := ParseDate(tc.in, now)
		switch {
		case tc.wantErr != nil:
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ParseDate(%q) err = %v, want %v", tc.in, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("ParseDate(%q) = %q, want %q (with future err)", tc.in, got, tc.want)
			}
		case tc.want == "":
			if err == nil {
				t.Errorf("ParseDate(%q) = %q, want error", tc.in, got)
			}
		default:
			if err != nil {
				t.Errorf("ParseDate(%q) unexpected error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("ParseDate(%q) = %q, want %q", tc.in, got, tc.want)
			}
		}
	}
}
