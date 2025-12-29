package html

type InvoiceView struct {
	Number     string
	IssuedAt   string
	Start      string
	End        string
	Consultant ConsultantView
	Company    CompanyView
	Contract   ContractView
	Activities []ActivityView
	Totals     TotalsView
}

type ConsultantView struct {
	FullName     string
	Country      string
	Zip          string
	Region       *string
	City         string
	AddressLine1 string
	AddressLine2 *string
	TaxNumber    string
	BankName     string
	BankAddress  string
	BankCountry  string
	BankIBAN     string
	BankBIC      string
}

type CompanyView struct {
	Name         string
	NameShort    *string
	TaxNumber    string
	Country      string
	Zip          string
	City         string
	Region       *string
	AddressLine1 string
	AddressLine2 *string
}

type ContractView struct {
	HourlyRateFormatted string
	Currency            string
	CurrencySign        string
	OrderNumber         string
	PaymentTerms        *string
}

type ActivityView struct {
	Title               string
	TotalHoursFormatted string
	HourlyRateFormatted string
	SubtotalFormatted   string
	Entries             []EntryView
}

type EntryView struct {
	DateFormatted  string
	TicketCode     string
	HoursFormatted string
}

type TotalsView struct {
	TotalHoursFormatted string
	SubtotalFormatted   string
}
