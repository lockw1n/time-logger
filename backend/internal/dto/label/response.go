package label

type Response struct {
	ID        uint64  `json:"id"`
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
}
