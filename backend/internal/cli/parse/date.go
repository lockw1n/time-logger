package parse

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// DateLayout is the canonical date format used across the CLI and the backend.
const DateLayout = "2006-01-02"

// ErrFutureDate is returned by ParseDate when the resolved date is after today.
// Logging future time is almost always a typo, so the command layer turns this
// into a warning-level confirm (skippable with --yes) rather than a hard error.
var ErrFutureDate = errors.New("date is in the future")

var weekdays = map[string]time.Weekday{
	"sun": time.Sunday,
	"mon": time.Monday,
	"tue": time.Tuesday,
	"wed": time.Wednesday,
	"thu": time.Thursday,
	"fri": time.Friday,
	"sat": time.Saturday,
}

// ParseDate normalizes user input to YYYY-MM-DD relative to now. It accepts an
// explicit YYYY-MM-DD, "today" (also the empty default), "yesterday", and the
// three-letter weekday names (most recent such weekday, today allowed). now is
// injected so "yesterday"/weekday resolution is testable.
//
// A resolved date after today is still returned (normalized) but paired with
// ErrFutureDate so the caller can confirm rather than silently reject.
func ParseDate(s string, now time.Time) (string, error) {
	// date-only anchor in now's location
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	s = strings.TrimSpace(strings.ToLower(s))

	var resolved time.Time
	switch {
	case s == "" || s == "today":
		return today.Format(DateLayout), nil
	case s == "yesterday":
		resolved = today.AddDate(0, 0, -1)
	default:
		if wd, ok := weekdays[s]; ok {
			resolved = mostRecentWeekday(today, wd)
		} else {
			d, err := time.ParseInLocation(DateLayout, s, now.Location())
			if err != nil {
				return "", fmt.Errorf("invalid date %q (accepted: YYYY-MM-DD, today, yesterday, mon-sun)", s)
			}
			resolved = d
		}
	}

	out := resolved.Format(DateLayout)
	if resolved.After(today) {
		return out, ErrFutureDate
	}
	return out, nil
}

// mostRecentWeekday returns the latest date <= today falling on wd (today itself
// qualifies when it is that weekday).
func mostRecentWeekday(today time.Time, wd time.Weekday) time.Time {
	delta := (int(today.Weekday()) - int(wd) + 7) % 7
	return today.AddDate(0, 0, -delta)
}
