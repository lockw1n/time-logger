package domain

type Consultant struct {
	ID           uint64
	FirstName    string
	MiddleName   *string
	LastName     string
	AddressLine1 string
	AddressLine2 *string
	Zip          string
	City         string
	Region       *string
	Country      string
	TaxNumber    string
	BankName     string
	BankAddress  string
	BankCountry  string
	BankIBAN     string
	BankBIC      string
}
