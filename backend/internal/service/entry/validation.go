package entry

import (
	"time"

	"github.com/lockw1n/time-logger/internal/constants"
)

func validateDurationMinutesRange(minutes int) error {
	if minutes <= 0 || minutes > 1440 {
		return ErrInvalidDurationMinutes
	}
	return nil
}

func validateDurationMinutesQuarter(minutes int) error {
	if minutes%15 != 0 {
		return ErrInvalidDurationMinutesQuarter
	}
	return nil
}

func validateDateFormat(date string) error {
	if _, err := time.Parse(constants.DateFormat, date); err != nil {
		return ErrInvalidDateFormat
	}
	return nil
}
