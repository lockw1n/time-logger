package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lockw1n/time-logger/internal/invoice"
)

type InvoiceHandler struct {
	Service        *EntryService
	CompanyService *CompanyService
	ConsultantSrv  *ConsultantService
}

func NewInvoiceHandler(service *EntryService, companyService *CompanyService, consultantService *ConsultantService) *InvoiceHandler {
	return &InvoiceHandler{
		Service:        service,
		CompanyService: companyService,
		ConsultantSrv:  consultantService,
	}
}

type invoiceRequest struct {
	Month       string `json:"month" binding:"required"` // YYYY-MM
	Currency    string `json:"currency"`                 // optional override (defaults to €)
	InvoiceDate string `json:"invoice_date"`             // optional display-friendly date, e.g. 02.01.2006
	PeriodStart string `json:"period_start"`             // optional override for displayed period
	PeriodEnd   string `json:"period_end"`               // optional override for displayed period
}

func displayDate(dateStr string) string {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return parsed.Format("02.01.2006")
}

func invoiceNumberForMonth(monthStr string) (string, error) {
	parsed, err := time.Parse("2006-01", strings.TrimSpace(monthStr))
	if err != nil {
		return "", err
	}
	nextMonth := parsed.AddDate(0, 1, 0)
	return nextMonth.Format("20060102"), nil
}

func joinName(parts ...string) string {
	trimmed := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			trimmed = append(trimmed, p)
		}
	}
	return strings.Join(trimmed, " ")
}

func (h *InvoiceHandler) GeneratePDF(c *gin.Context) {
	var req invoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	company, err := h.CompanyService.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch company"})
		return
	}
	if company == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company not configured"})
		return
	}

	consultant, err := h.ConsultantSrv.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch consultant"})
		return
	}
	if consultant == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "consultant not configured"})
		return
	}

	invoiceNumber, err := invoiceNumberForMonth(req.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month is required and must be in YYYY-MM format"})
		return
	}

	report, err := h.Service.MonthlySummary(req.Month)
	if err != nil {
		if err == ErrBadMonth {
			c.JSON(http.StatusBadRequest, gin.H{"error": "month is required and must be in YYYY-MM format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build monthly summary"})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "€"
	}

	paymentCondition := company.Payment
	if paymentCondition == "" {
		paymentCondition = "Net 14 days"
	}

	invoiceDate := req.InvoiceDate
	if invoiceDate == "" {
		invoiceDate = time.Now().Format("02.01.2006")
	}

	periodStart := req.PeriodStart
	periodEnd := req.PeriodEnd
	if periodStart == "" {
		periodStart = displayDate(report.Start)
	}
	if periodEnd == "" {
		periodEnd = displayDate(report.End)
	}

	items := make([]invoice.Item, 0, len(report.Items))
	totalAmount := 0.0
	for _, it := range report.Items {
		rate := consultant.HourlyRate
		amount := it.TotalHours * rate
		items = append(items, invoice.Item{
			Description: it.Label,
			Hours:       it.TotalHours,
			UnitPrice:   rate,
			Amount:      amount,
		})
		totalAmount += amount
	}

	fullName := joinName(consultant.FirstName, consultant.MiddleName, consultant.LastName)

	data := invoice.Data{
		ConsultantName:       fullName,
		ConsultantFirstName:  consultant.FirstName,
		ConsultantMiddleName: consultant.MiddleName,
		ConsultantLastName:   consultant.LastName,
		ConsultantAddress:    consultant.Address1,
		ConsultantAddress2:   consultant.Address2,
		ConsultantTaxNumber:  consultant.TaxNumber,
		ConsultantRegion:     consultant.Region,
		ConsultantZip:        consultant.Zip,
		ConsultantCity:       consultant.City,
		ConsultantCountry:    consultant.Country,
		CompanyName:          company.Name,
		CompanyUID:           company.UID,
		CompanyStreet:        company.AddressLine1,
		CompanyCity:          company.City,
		CompanyCountry:       company.Country,
		CompanyZip:           company.Zip,
		OrderNumber:          consultant.OrderNumber,
		InvoiceNumber:        invoiceNumber,
		InvoiceDate:          invoiceDate,
		PeriodStart:          periodStart,
		PeriodEnd:            periodEnd,
		HourlyRate:           consultant.HourlyRate,
		Currency:             currency,
		PaymentCondition:     paymentCondition,
		BankName:             consultant.BankName,
		BankAddress:          consultant.BankAddress,
		IBAN:                 consultant.IBAN,
		BIC:                  consultant.BIC,
		BankCountry:          consultant.BankCountry,
		Items:                items,
		TotalHours:           report.TotalHours,
		TotalAmount:          totalAmount,
	}

	html, err := invoice.RenderHTML(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render invoice HTML"})
		return
	}

	pdf, err := invoice.GeneratePDF(html)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("wkhtmltopdf error: %v", err)})
		return
	}

	filename := fmt.Sprintf("invoice-%s.pdf", report.Month)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdf)
}
