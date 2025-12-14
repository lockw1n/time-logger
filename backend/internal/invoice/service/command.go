package service

type GenerateMonthlyInvoiceCommand struct {
	Month        string // YYYY-MM
	ConsultantID uint64
	CompanyID    uint64
}
