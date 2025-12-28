package render

import (
	"fmt"
	"time"

	"github.com/lockw1n/time-logger/internal/constants"
	consultantdomain "github.com/lockw1n/time-logger/internal/consultant/domain"
	"github.com/lockw1n/time-logger/internal/invoice/domain"
	"github.com/lockw1n/time-logger/internal/invoice/render/excel"
	"github.com/lockw1n/time-logger/internal/invoice/render/html"
)

func MapInvoiceToHTMLView(invoice domain.Invoice) html.InvoiceView {
	return html.InvoiceView{
		Number:     invoice.Number,
		IssuedAt:   formatDate(invoice.IssuedAt),
		Start:      formatDate(invoice.Start),
		End:        formatDate(invoice.End),
		Consultant: mapConsultantToHTMLView(invoice),
		Company:    mapCompanyToHTMLView(invoice),
		Contract:   mapContractToHTMLView(invoice),
		Activities: mapActivitiesToHTMLView(invoice.Activities, invoice.Company.NameShort),
		Totals:     mapTotalsToHTMLView(invoice.Totals),
	}
}

func mapConsultantToHTMLView(invoice domain.Invoice) html.ConsultantView {
	consultant := invoice.Consultant

	return html.ConsultantView{
		FullName:     ConsultantFullName(consultant),
		Country:      consultant.Country,
		Zip:          consultant.Zip,
		Region:       stringPtr(consultant.Region),
		City:         consultant.City,
		AddressLine1: consultant.AddressLine1,
		AddressLine2: stringPtr(consultant.AddressLine2),
		TaxNumber:    consultant.TaxNumber,
		BankName:     consultant.BankName,
		BankAddress:  consultant.BankAddress,
		BankCountry:  consultant.BankCountry,
		BankIBAN:     consultant.BankIBAN,
		BankBIC:      consultant.BankBIC,
	}
}

func mapCompanyToHTMLView(invoice domain.Invoice) html.CompanyView {
	company := invoice.Company

	return html.CompanyView{
		Name:         company.Name,
		NameShort:    stringPtr(company.NameShort),
		TaxNumber:    company.TaxNumber,
		Country:      company.Country,
		Zip:          company.Zip,
		City:         company.City,
		Region:       stringPtr(company.Region),
		AddressLine1: company.AddressLine1,
		AddressLine2: stringPtr(company.AddressLine2),
	}
}

func mapContractToHTMLView(invoice domain.Invoice) html.ContractView {
	contract := invoice.Contract

	return html.ContractView{
		OrderNumber:         contract.OrderNumber,
		PaymentTerms:        stringPtr(contract.PaymentTerms),
		HourlyRateFormatted: formatMoney(contract.HourlyRate, contract.Currency),
		Currency:            contract.Currency,
	}
}

func mapActivitiesToHTMLView(activities []domain.InvoiceActivity, companyShort *string) []html.ActivityView {
	out := make([]html.ActivityView, 0, len(activities))

	for _, activity := range activities {
		out = append(out, html.ActivityView{
			Title:               activityTitle(companyShort, activity.Name),
			TotalHoursFormatted: formatHours(activity.TotalHours),
			HourlyRateFormatted: formatMoney(activity.HourlyRate, ""),
			SubtotalFormatted:   formatCents(activity.Subtotal),
			Entries:             mapEntriesToHTMLView(activity.Entries),
		})
	}

	return out
}

func mapEntriesToHTMLView(entries []domain.InvoiceEntry) []html.EntryView {
	out := make([]html.EntryView, 0, len(entries))

	for _, entry := range entries {
		out = append(out, html.EntryView{
			DateFormatted:  formatDate(entry.Date),
			TicketCode:     entry.TicketCode,
			HoursFormatted: formatHours(entry.Hours),
		})
	}

	return out
}

func mapTotalsToHTMLView(totals domain.InvoiceTotals) html.TotalsView {
	return html.TotalsView{
		TotalHoursFormatted: formatHours(totals.TotalHours),
		SubtotalFormatted:   formatCents(totals.Subtotal),
	}
}

func MapInvoiceToExcelView(invoice domain.Invoice) excel.InvoiceView {
	return excel.InvoiceView{
		Number:     invoice.Number,
		Activities: mapActivitiesToExcelView(invoice.Activities, invoice.Company.NameShort),
	}
}

func mapActivitiesToExcelView(activities []domain.InvoiceActivity, companyShort *string) []excel.ActivityView {
	out := make([]excel.ActivityView, 0, len(activities))

	for _, activity := range activities {
		out = append(out, excel.ActivityView{
			Title:   activityTitle(companyShort, activity.Name),
			Entries: mapEntriesToExcelView(activity.Entries),
		})
	}

	return out
}

func mapEntriesToExcelView(entries []domain.InvoiceEntry) []excel.EntryView {
	out := make([]excel.EntryView, 0, len(entries))

	for _, entry := range entries {
		out = append(out, excel.EntryView{
			DateFormatted: formatDate(entry.Date),
			TicketCode:    entry.TicketCode,
			Hours:         entry.Hours,
		})
	}

	return out
}

func ConsultantFullName(consultant consultantdomain.Consultant) string {
	if mid := stringPtr(consultant.MiddleName); mid != nil {
		return consultant.FirstName + " " + *mid + " " + consultant.LastName
	}
	return consultant.FirstName + " " + consultant.LastName
}

func activityTitle(companyShort *string, activityName string) string {
	if companyShort != nil && *companyShort != "" {
		return *companyShort + " - " + activityName
	}
	return activityName
}

func formatDate(t time.Time) string {
	return t.Format(constants.ResponseDateFormat)
}

func formatHours(h float64) string {
	return fmt.Sprintf("%.2f", h)
}

func formatMoney(amount float64, currency string) string {
	if currency == "" {
		return fmt.Sprintf("%.2f", amount)
	}
	sign := currencySign(currency)
	return fmt.Sprintf("%s%.2f", sign, amount)
}

func currencySign(currency string) string {
	switch currency {
	case "EUR":
		return "€"
	case "USD":
		return "$"
	case "UAH":
		return "₴"
	default:
		return currency
	}
}

func formatCents(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}

func stringPtr(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
