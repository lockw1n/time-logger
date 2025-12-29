package service

import "strings"

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i CreateEntryInput) Validate() error {
	if i.ConsultantID == 0 {
		return validationError("consultant_id is required")
	}
	if i.CompanyID == 0 {
		return validationError("company_id is required")
	}
	if strings.TrimSpace(i.TicketCode) == "" {
		return validationError("ticket code is required and should not be empty")
	}
	if i.ActivityID == 0 {
		return validationError("activity_id is required")
	}
	if i.Date.IsZero() {
		return validationError("date is required")
	}
	if i.DurationMinutes <= 0 || i.DurationMinutes > 1440 {
		return validationError("duration_minutes should be between 1 and 1440")
	}

	return nil
}

func (i UpdateEntryInput) Validate() error {
	if i.DurationMinutes != nil && (*i.DurationMinutes <= 0 || *i.DurationMinutes > 1440) {
		return validationError("duration_minutes should be between 1 and 1440")
	}

	return nil
}
