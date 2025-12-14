package render

type Company struct {
	Name      string `json:"name"`
	NameShort string `json:"nameShort"`
	TaxNumber string `json:"taxNumber"`

	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Zip          string `json:"zip"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`

	PaymentTerms string `json:"paymentTerms"`
}
