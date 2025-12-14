package render

type Group struct {
	Label string `json:"label"`

	TotalHours          float64 `json:"totalHours"`
	TotalHoursFormatted string  `json:"totalHoursFormatted"`

	Subtotal          int64  `json:"subtotal"`
	SubtotalFormatted string `json:"subtotalFormatted"`

	Rows []Row `json:"rows"`
}
