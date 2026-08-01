package parse

import (
	"fmt"
	"strings"
	"time"
)

// monthLayout is the YYYY-MM form accepted by `tl invoice --month`.
const monthLayout = "2006-01"

// MonthRange expands a "YYYY-MM" month into its first and last calendar day as
// YYYY-MM-DD. Month length (including February in leap and non-leap years) is
// handled by date arithmetic: last = firstOfMonth + 1 month - 1 day.
func MonthRange(s string) (start, end string, err error) {
	s = strings.TrimSpace(s)
	t, err := time.Parse(monthLayout, s)
	if err != nil {
		return "", "", fmt.Errorf("invalid month %q (accepted: YYYY-MM)", s)
	}
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1)
	return first.Format(DateLayout), last.Format(DateLayout), nil
}

// PreviousMonth returns the "YYYY-MM" of the calendar month before now — the
// default period for `tl invoice`, since invoices are generated after a period
// closes. December→January rolls the year back correctly.
func PreviousMonth(now time.Time) string {
	first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return first.AddDate(0, -1, 0).Format(monthLayout)
}
