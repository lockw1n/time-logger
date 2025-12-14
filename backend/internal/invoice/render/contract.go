package render

type Contract struct {
	OrderNumber         string  `json:"orderNumber"`
	HourlyRate          float64 `json:"hourlyRate"`
	HourlyRateFormatted string  `json:"hourlyRateFormatted"`
}
