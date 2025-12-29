package domain

type Company struct {
	ID           uint64
	Name         string
	NameShort    *string
	TaxNumber    string
	AddressLine1 string
	AddressLine2 *string
	Zip          string
	City         string
	Region       *string
	Country      string
}
