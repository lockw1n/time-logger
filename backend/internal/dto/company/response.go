package company

type Response struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name"`
	NameShort    *string `json:"name_short"`
	TaxNumber    *string `json:"tax_number"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Zip          string  `json:"zip"`
	City         string  `json:"city"`
	Region       *string `json:"region"`
	Country      string  `json:"country"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
