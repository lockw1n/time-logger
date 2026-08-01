package parse

import "time"

// WeekBounds returns the Monday and Sunday (inclusive) of the ISO week
// containing the given YYYY-MM-DD date, both as YYYY-MM-DD strings.
func WeekBounds(date string) (monday, sunday string, err error) {
	d, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", "", err
	}

	// Go weekdays run Sunday=0..Saturday=6; shift so Monday is the anchor.
	offset := (int(d.Weekday()) + 6) % 7
	mon := d.AddDate(0, 0, -offset)
	sun := mon.AddDate(0, 0, 6)

	return mon.Format(DateLayout), sun.Format(DateLayout), nil
}

// WeekDays returns the seven YYYY-MM-DD dates Monday..Sunday for the week
// containing date. Used by the week renderer to index per-day columns.
func WeekDays(date string) ([]string, error) {
	monday, _, err := WeekBounds(date)
	if err != nil {
		return nil, err
	}

	mon, _ := time.Parse(DateLayout, monday)
	days := make([]string, 7)
	for i := range days {
		days[i] = mon.AddDate(0, 0, i).Format(DateLayout)
	}
	return days, nil
}
