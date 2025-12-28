package domain

type Ticket struct {
	ID          uint64
	CompanyID   uint64
	Code        string
	Title       *string
	Label       *string
	Description *string
}
