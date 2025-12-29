package domain

import "time"

type Entry struct {
	ID              uint64
	ConsultantID    uint64
	CompanyID       uint64
	ContractID      uint64
	TicketID        uint64
	ActivityID      uint64
	Date            time.Time
	DurationMinutes int
	Comment         *string
}
