package service

type CreateTicketInput struct {
	CompanyID   uint64
	Code        string
	Title       *string
	Label       *string
	Description *string
}

type UpdateTicketInput struct {
	Title       *string
	Label       *string
	Description *string
}
