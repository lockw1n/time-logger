package render

type Group struct {
	Label string `json:"label"`

	TotalHours          float64 `json:"totalHours"`
	TotalHoursFormatted string  `json:"totalHoursFormatted"`

	HourlyRate          float64 `json:"hourlyRate"`
	HourlyRateFormatted string  `json:"hourlyRateFormatted"`

	Subtotal          int64  `json:"subtotal"`
	SubtotalFormatted string `json:"subtotalFormatted"`

	Rows []Row `json:"rows"`
}
