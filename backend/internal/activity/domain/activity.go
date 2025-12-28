package domain

type Activity struct {
	ID        uint64
	CompanyID uint64
	Name      string
	Color     *string
	Billable  bool
	Priority  int
}
