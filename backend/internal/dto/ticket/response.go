package ticket

type Response struct {
	ID          uint64  `json:"id"`
	CompanyID   uint64  `json:"company_id"`
	Code        string  `json:"code"`
	Label       string  `json:"label"`
	Description *string `json:"description"`
}
