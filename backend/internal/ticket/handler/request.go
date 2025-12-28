package handler

type CreateTicketRequest struct {
	CompanyID   uint64  `json:"company_id"`
	Code        string  `json:"code"`
	Title       *string `json:"title"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

type UpdateTicketRequest struct {
	Title       *string `json:"title"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}
