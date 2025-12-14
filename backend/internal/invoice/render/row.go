package render

type Row struct {
	Date          string `json:"date"`
	DateFormatted string `json:"dateFormatted"`

	TicketCode  string `json:"ticketCode"`
	Description string `json:"description"`

	Hours          float64 `json:"hours"`
	HoursFormatted string  `json:"hoursFormatted"`

	Amount          int64  `json:"amount"`
	AmountFormatted string `json:"amountFormatted"`
}
