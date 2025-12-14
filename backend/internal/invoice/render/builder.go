package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/lockw1n/time-logger/internal/invoice/domain"
)

func BuildInvoice(invoice domain.Invoice) Invoice {
	return Invoice{
		Invoice:    buildInvoiceMeta(invoice),
		Period:     buildPeriod(invoice.Period),
		Consultant: buildConsultant(invoice.Consultant),
		Company:    buildCompany(invoice.Company),
		Contract:   buildContract(invoice.Contract, invoice.Currency),
		Groups:     buildGroups(invoice.Groups, invoice.Currency),
		Totals:     buildTotals(invoice.Totals, invoice.Currency),
	}
}

func buildInvoiceMeta(invoice domain.Invoice) InvoiceMeta {
	return InvoiceMeta{
		Number: invoice.Number,

		IssuedAt:          invoice.IssuedAt.Format("2006-01-02"),
		IssuedAtFormatted: invoice.IssuedAt.Format("02.01.2006"),

		DueAt:          invoice.DueAt.Format("2006-01-02"),
		DueAtFormatted: invoice.DueAt.Format("02.01.2006"),

		Currency:       invoice.Currency,
		CurrencySymbol: currencySymbol(invoice.Currency),
	}
}

func currencySymbol(currency string) string {
	switch currency {
	case "EUR":
		return "€"
	case "USD":
		return "$"
	case "GBP":
		return "£"
	default:
		return currency
	}
}

func buildPeriod(period domain.Period) Period {
	t, err := time.Parse("2006-01", period.Month)
	if err != nil {
		t = period.Start
	}

	return Period{
		Month: period.Month,
		Label: t.Format("January 2006"),

		Start: period.Start.Format("02.01.2006"),
		End:   period.End.Format("02.01.2006"),
	}
}

func buildConsultant(consultant domain.Consultant) Consultant {
	return Consultant{
		FullName: buildFullName(
			consultant.FirstName,
			consultant.MiddleName,
			consultant.LastName,
		),

		AddressLine1: consultant.AddressLine1,
		AddressLine2: consultant.AddressLine2,
		Zip:          consultant.Zip,
		City:         consultant.City,
		Region:       consultant.Region,
		Country:      consultant.Country,

		TaxNumber: consultant.TaxNumber,

		BankName:    consultant.BankName,
		BankAddress: consultant.BankAddress,
		BankCountry: consultant.BankCountry,
		BankIban:    consultant.BankIban,
		BankBic:     consultant.BankBic,
	}
}

func buildFullName(first, middle, last string) string {
	if middle == "" {
		return first + " " + last
	}
	return first + " " + middle + " " + last
}

func buildCompany(company domain.Company) Company {
	return Company{
		Name:      company.Name,
		NameShort: company.NameShort,
		TaxNumber: company.TaxNumber,

		AddressLine1: company.AddressLine1,
		AddressLine2: company.AddressLine2,
		Zip:          company.Zip,
		City:         company.City,
		Region:       company.Region,
		Country:      company.Country,

		PaymentTerms: company.PaymentTerms,
	}
}

func buildContract(contract domain.Contract, currency string) Contract {
	return Contract{
		OrderNumber: contract.OrderNumber,
		HourlyRate:  contract.HourlyRate,
		HourlyRateFormatted: formatHourlyRate(
			contract.HourlyRate,
			currency,
		),
	}
}

func formatHourlyRate(rate float64, currency string) string {
	symbol := currencySymbol(currency)

	// Convert decimal separator: 10.00 -> 10,00
	value := fmt.Sprintf("%.2f", rate)
	value = strings.Replace(value, ".", ",", 1)

	return symbol + value
}

func buildGroups(groups []domain.Group, currency string) []Group {
	result := make([]Group, 0, len(groups))

	for _, group := range groups {
		result = append(result, buildGroup(group, currency))
	}

	return result
}

func buildGroup(group domain.Group, currency string) Group {
	rows := make([]Row, 0, len(group.Rows))
	for _, row := range group.Rows {
		rows = append(rows, buildRow(row, currency))
	}

	return Group{
		Label: group.Label,

		TotalHours:          group.TotalHours,
		TotalHoursFormatted: formatHours(group.TotalHours),

		HourlyRate:          group.HourlyRate,
		HourlyRateFormatted: formatHourlyRate(group.HourlyRate, currency),

		Subtotal:          group.Subtotal,
		SubtotalFormatted: formatMoney(group.Subtotal, currency),

		Rows: rows,
	}
}

func buildRow(row domain.Row, currency string) Row {
	return Row{
		Date:          row.Date.Format("2006-01-02"),
		DateFormatted: row.Date.Format("02.01.2006"),

		TicketCode:  row.TicketCode,
		Description: row.Description,

		Hours:          row.Hours,
		HoursFormatted: formatHours(row.Hours),

		Amount:          row.Amount,
		AmountFormatted: formatMoney(row.Amount, currency),
	}
}

func buildTotals(totals domain.Totals, currency string) Totals {
	return Totals{
		TotalHours: totals.TotalHours,
		TotalHoursFormatted: formatHours(
			totals.TotalHours,
		),

		Subtotal: totals.Subtotal,
		SubtotalFormatted: formatMoney(
			totals.Subtotal,
			currency,
		),
	}
}

func formatHours(hours float64) string {
	return fmt.Sprintf("%.2f", hours)
}

func formatMoney(amount int64, currency string) string {
	symbol := currencySymbol(currency)

	value := float64(amount) / 100.0
	str := fmt.Sprintf("%.2f", value)
	str = strings.Replace(str, ".", ",", 1)

	return symbol + str
}
