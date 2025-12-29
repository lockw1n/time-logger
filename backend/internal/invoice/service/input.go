package service

import "time"

type GenerateInvoiceInput struct {
	ConsultantID uint64
	CompanyID    uint64
	Start        time.Time
	End          time.Time
}
