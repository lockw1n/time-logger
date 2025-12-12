package label

type Request struct {
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
}
