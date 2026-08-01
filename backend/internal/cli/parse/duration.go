// Package parse turns human-typed durations and dates into the integer-minute
// and YYYY-MM-DD forms the backend expects. It is pure (no I/O) so every
// accepted and rejected form is table-testable; the CLI parses at the edge and
// works in minutes internally.
package parse

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// MaxMinutes is the largest duration the CLI (and the backend) accepts for a
// single entry: 24h. Exported so callers that build a total themselves — the
// merge flow — can enforce the same ceiling.
const MaxMinutes = 1440

var errDurationFormat = errors.New(
	"invalid duration (accepted: 90m, 90, 1h, 1h30m, 1h30, 1.5h, 0.25h)")

var (
	// fractionalHour matches the fractional-hour form "1.5h" / "0.25h";
	// group 1 is the (possibly decimal) hour count. Fractions need the h suffix.
	fractionalHour = regexp.MustCompile(`^(\d+(?:\.\d+)?)h$`)
	// hoursMinutes matches "1h", "1h30m", "1h30", "90m", "90";
	// group 1 is whole hours, group 2 is minutes (both optional).
	hoursMinutes = regexp.MustCompile(`^(?:(\d+)h)?(?:(\d+)m?)?$`)
)

// ParseDuration parses a human duration into whole minutes. The result must be
// > 0 and <= 1440. Fractional values are allowed only with the "h" suffix and
// are rounded to the nearest minute.
func ParseDuration(s string) (int, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, errDurationFormat
	}

	if m := fractionalHour.FindStringSubmatch(s); m != nil {
		hours, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return 0, errDurationFormat
		}
		return validateMinutes(int(math.Round(hours * 60)))
	}

	m := hoursMinutes.FindStringSubmatch(s)
	// The regex matches the empty string and a lone "h"; require at least one digit.
	if m == nil || (m[1] == "" && m[2] == "") {
		return 0, errDurationFormat
	}

	total := 0
	if m[1] != "" {
		h, err := strconv.Atoi(m[1])
		if err != nil {
			return 0, errDurationFormat
		}
		total += h * 60
	}
	if m[2] != "" {
		mins, err := strconv.Atoi(m[2])
		if err != nil {
			return 0, errDurationFormat
		}
		total += mins
	}

	return validateMinutes(total)
}

func validateMinutes(m int) (int, error) {
	if m <= 0 || m > MaxMinutes {
		return 0, fmt.Errorf("duration must be between 1m and 24h, got %s", FormatMinutes(m))
	}
	return m, nil
}

// FormatMinutes renders whole minutes as "7h30m", "45m", "8h". Used by all
// human-facing rendering.
func FormatMinutes(m int) string {
	if m <= 0 {
		return "0m"
	}

	h, mins := m/60, m%60
	switch {
	case h == 0:
		return fmt.Sprintf("%dm", mins)
	case mins == 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dh%dm", h, mins)
	}
}
