package handler

type GenerateInvoiceRequest struct {
	ConsultantID uint64 `form:"consultant_id" binding:"required"`
	CompanyID    uint64 `form:"company_id" binding:"required"`
	Start        string `form:"start" binding:"required"`
	End          string `form:"end" binding:"required"`
}
