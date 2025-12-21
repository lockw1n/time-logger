package service

import (
	"time"

	"github.com/lockw1n/time-logger/internal/constants"
)

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(constants.DateFormat, date)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return t, nil
}

func validateReportScope(consultantID uint64, companyID uint64) error {
	if consultantID == 0 {
		return ErrMissingConsultantID
	}
	if companyID == 0 {
		return ErrMissingCompanyID
	}
	return nil
}

func validateDateRange(start time.Time, end time.Time) error {
	if end.Before(start) {
		return ErrInvalidDateRange
	}
	return nil
}
