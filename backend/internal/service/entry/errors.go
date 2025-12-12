package entry

import "errors"

var (
	ErrNotFound                      = errors.New("entry not found")
	ErrInvalidDurationMinutes        = errors.New("duration in minutes must be in between 15 and 1440")
	ErrInvalidDurationMinutesQuarter = errors.New("duration must be in 15-minute increments")
	ErrInvalidDateFormat             = errors.New("date must be in YYYY-MM-DD format")
)
