package service

import "time"

type CreateEntryInput struct {
	CompanyID       uint64
	ActivityID      uint64
	TicketCode      string
	Date            time.Time
	DurationMinutes int
	Comment         *string
}

type UpdateEntryInput struct {
	DurationMinutes *int
	Comment         *string
}
