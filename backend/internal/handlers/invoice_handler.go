package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lockw1n/time-logger/internal/invoice"
)

type InvoiceHandler struct {
	Service *EntryService
}

func NewInvoiceHandler(service *EntryService) *InvoiceHandler {
	return &InvoiceHandler{Service: service}
}

type invoiceRequest struct {
	Month             string  `json:"month" binding:"required"` // YYYY-MM
	InvoiceNumber     string  `json:"invoice_number"`
	OrderNumber       string  `json:"order_number"`
	HourlyRate        float64 `json:"hourly_rate"`
	Currency          string  `json:"currency"`
	InvoiceDate       string  `json:"invoice_date"` // display-friendly date, e.g. 02.01.2006
	PeriodStart       string  `json:"period_start"`
	PeriodEnd         string  `json:"period_end"`
	ConsultantName    string  `json:"consultant_name"`
	ConsultantAddress string  `json:"consultant_address"`
	ConsultantTax     string  `json:"consultant_tax_number"`
	CompanyName       string  `json:"company_name"`
	CompanyUID        string  `json:"company_uid"`
	CompanyStreet     string  `json:"company_street"`
	CompanyCity       string  `json:"company_city"`
	CompanyCountry    string  `json:"company_country"`
	PaymentCondition  string  `json:"payment_condition"`
	BankName          string  `json:"bank_name"`
	BankAddress       string  `json:"bank_address"`
	IBAN              string  `json:"iban"`
	BIC               string  `json:"bic"`
	BankCountry       string  `json:"bank_country"`
}

func displayDate(dateStr string) string {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return parsed.Format("02.01.2006")
}

func (h *InvoiceHandler) GeneratePDF(c *gin.Context) {
	var req invoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.HourlyRate <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hourly_rate must be greater than zero"})
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

	paymentCondition := req.PaymentCondition
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
		amount := it.TotalHours * req.HourlyRate
		items = append(items, invoice.Item{
			Description: it.Label,
			Hours:       it.TotalHours,
			UnitPrice:   req.HourlyRate,
			Amount:      amount,
		})
		totalAmount += amount
	}

	data := invoice.Data{
		ConsultantName:      req.ConsultantName,
		ConsultantAddress:   req.ConsultantAddress,
		ConsultantTaxNumber: req.ConsultantTax,
		CompanyName:         req.CompanyName,
		CompanyUID:          req.CompanyUID,
		CompanyStreet:       req.CompanyStreet,
		CompanyCity:         req.CompanyCity,
		CompanyCountry:      req.CompanyCountry,
		OrderNumber:         req.OrderNumber,
		InvoiceNumber:       req.InvoiceNumber,
		InvoiceDate:         invoiceDate,
		PeriodStart:         periodStart,
		PeriodEnd:           periodEnd,
		HourlyRate:          req.HourlyRate,
		Currency:            currency,
		PaymentCondition:    paymentCondition,
		BankName:            req.BankName,
		BankAddress:         req.BankAddress,
		IBAN:                req.IBAN,
		BIC:                 req.BIC,
		BankCountry:         req.BankCountry,
		Items:               items,
		TotalHours:          report.TotalHours,
		TotalAmount:         totalAmount,
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
