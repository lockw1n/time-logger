package domain

import (
	"time"

	companydomain "github.com/lockw1n/time-logger/internal/company/domain"
	consultantdomain "github.com/lockw1n/time-logger/internal/consultant/domain"
	contractdomain "github.com/lockw1n/time-logger/internal/contract/domain"
)

type Invoice struct {
	Number     string
	IssuedAt   time.Time
	Start      time.Time
	End        time.Time
	Consultant consultantdomain.Consultant
	Company    companydomain.Company
	Contract   contractdomain.Contract
	Activities []InvoiceActivity
	Totals     InvoiceTotals
}

type InvoiceActivity struct {
	Name       string
	Entries    []InvoiceEntry
	TotalHours float64
	HourlyRate float64
	Subtotal   int64 // smallest currency unit (e.g. cents)
}

type InvoiceEntry struct {
	Date       time.Time
	TicketCode string
	Hours      float64
}

type InvoiceTotals struct {
	TotalHours float64
	Subtotal   int64 // smallest currency unit
}
