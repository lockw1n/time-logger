package ticket

type Request struct {
	CompanyID   uint64  `json:"company_id"`
	Code        string  `json:"code"`
	Label       string  `json:"label"`
	Description *string `json:"description"`
}
