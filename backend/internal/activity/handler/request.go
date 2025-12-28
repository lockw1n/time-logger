package handler

type CreateActivityRequest struct {
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
	Billable  bool    `json:"billable"`
	Priority  int     `json:"priority"`
}

type UpdateActivityRequest struct {
	Name     *string `json:"name"`
	Color    *string `json:"color"`
	Billable *bool   `json:"billable"`
	Priority *int    `json:"priority"`
}
