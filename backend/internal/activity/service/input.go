package service

type CreateActivityInput struct {
	CompanyID uint64
	Name      string
	Color     *string
	Billable  bool
	Priority  int
}

type UpdateActivityInput struct {
	Name     *string
	Color    *string
	Billable *bool
	Priority *int
}
