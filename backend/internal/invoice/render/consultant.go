package render

type Consultant struct {
	FullName string `json:"fullName"`

	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Zip          string `json:"zip"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`

	TaxNumber string `json:"taxNumber"`

	BankName    string `json:"bankName"`
	BankAddress string `json:"bankAddress"`
	BankCountry string `json:"bankCountry"`
	BankIban    string `json:"bankIban"`
	BankBic     string `json:"bankBic"`
}
