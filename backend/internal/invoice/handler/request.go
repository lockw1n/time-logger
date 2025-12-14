package handler

type GenerateMonthlyInvoiceRequest struct {
	Month        string `json:"month" binding:"required"` // "2025-12"
	ConsultantID uint64 `json:"consultant_id" binding:"required"`
	CompanyID    uint64 `json:"company_id" binding:"required"`
}
