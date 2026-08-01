package parse

import "testing"

func TestParseDuration(t *testing.T) {
	valid := []struct {
		in   string
		want int
	}{
		{"90m", 90},
		{"90", 90},
		{"1h", 60},
		{"1h30m", 90},
		{"1h30", 90},
		{"1.5h", 90},
		{"0.25h", 15},
		{"8h", 480},
		{"24h", 1440},
		{" 45M ", 45}, // trimmed + case-insensitive
		{"0.1h", 6},   // 6.0 -> 6
		{"0.08h", 5},  // 4.8 -> rounds to 5
	}
	for _, tc := range valid {
		got, err := ParseDuration(tc.in)
		if err != nil {
			t.Errorf("ParseDuration(%q) unexpected error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseDuration(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}

	invalid := []string{
		"", "0", "0m", "h", "m", "1.5", "1,5h", "abc",
		"-5m", "25h", "1441", "1h1441m", "1.5", "1h.5", "1.5m",
	}
	for _, in := range invalid {
		if got, err := ParseDuration(in); err == nil {
			t.Errorf("ParseDuration(%q) = %d, want error", in, got)
		}
	}
}

func TestFormatMinutes(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{450, "7h30m"},
		{45, "45m"},
		{480, "8h"},
		{0, "0m"},
		{-10, "0m"},
		{1, "1m"},
		{60, "1h"},
		{1440, "24h"},
	}
	for _, tc := range cases {
		if got := FormatMinutes(tc.in); got != tc.want {
			t.Errorf("FormatMinutes(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
