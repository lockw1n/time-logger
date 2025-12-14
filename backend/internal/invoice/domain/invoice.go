package domain

import "time"

type Invoice struct {
	Number   string
	IssuedAt time.Time
	DueAt    time.Time
	Currency string

	Period     Period
	Consultant Consultant
	Company    Company
	Contract   Contract

	Groups []Group
	Totals Totals
}
