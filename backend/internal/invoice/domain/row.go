package domain

import "time"

type Row struct {
	Date        time.Time
	TicketCode  string
	Description string

	Hours  float64
	Amount int64 // smallest currency unit
}
