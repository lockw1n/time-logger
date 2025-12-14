package service

import "errors"

var (
	ErrInvalidMonthFormat  = errors.New("invalid month format, expected YYYY-MM")
	ErrAssignmentNotFound  = errors.New("assignment not found for given month")
	ErrMultipleAssignments = errors.New("multiple assignments found for given month")
	ErrNoEntriesForPeriod  = errors.New("no entries found for given month")
)
