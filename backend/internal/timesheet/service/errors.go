package service

import "errors"

var (
	ErrInvalidDateFormat   = errors.New("date must be in YYYY-MM-DD format")
	ErrInvalidDateRange    = errors.New("start date must be before end date")
	ErrMissingConsultantID = errors.New("consultant_id is required")
	ErrMissingCompanyID    = errors.New("company_id is required")
)
