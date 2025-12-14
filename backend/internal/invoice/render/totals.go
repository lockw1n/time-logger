package render

type Totals struct {
	TotalHours          float64 `json:"totalHours"`
	TotalHoursFormatted string  `json:"totalHoursFormatted"`

	Subtotal          int64  `json:"subtotal"`
	SubtotalFormatted string `json:"subtotalFormatted"`
}
