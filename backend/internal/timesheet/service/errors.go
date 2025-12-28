package service

import "errors"

var (
	ErrTimesheetInvalid  = errors.New("invalid timesheet data")
	ErrTimesheetConflict = errors.New("timesheet conflict")
)
