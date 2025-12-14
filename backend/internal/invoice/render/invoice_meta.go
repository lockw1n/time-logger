package render

type InvoiceMeta struct {
	Number            string `json:"number"`
	IssuedAt          string `json:"issuedAt"`          // ISO
	IssuedAtFormatted string `json:"issuedAtFormatted"` // 31.12.2025
	DueAt             string `json:"dueAt"`
	DueAtFormatted    string `json:"dueAtFormatted"`
	Currency          string `json:"currency"`
	CurrencySymbol    string `json:"currencySymbol"`
}
