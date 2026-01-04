package handler

type GenerateInvoiceRequest struct {
	CompanyID uint64 `form:"company_id" binding:"required"`
	Start     string `form:"start" binding:"required"`
	End       string `form:"end" binding:"required"`
}
