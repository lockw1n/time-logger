package consultant

type Response struct {
	ID         uint64  `json:"id"`
	FirstName  string  `json:"first_name"`
	MiddleName *string `json:"middle_name"`
	LastName   string  `json:"last_name"`

	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Zip          string  `json:"zip"`
	City         string  `json:"city"`
	Region       *string `json:"region"`
	Country      string  `json:"country"`

	TaxNumber *string `json:"tax_number"`

	BankName    string `json:"bank_name"`
	BankAddress string `json:"bank_address"`
	BankCountry string `json:"bank_country"`
	BankIBAN    string `json:"bank_iban"`
	BankBIC     string `json:"bank_bic"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
