package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/excel"
	"github.com/lockw1n/time-logger/internal/invoice/domain"
	"github.com/lockw1n/time-logger/internal/invoice/render"
	"github.com/lockw1n/time-logger/internal/invoice/service"
	"github.com/lockw1n/time-logger/internal/pdf"
)

type Invoice struct {
	generator     service.InvoiceGenerator
	pdfRenderer   pdf.Renderer
	excelRenderer *excel.Renderer
}

func NewInvoice(
	generator service.InvoiceGenerator,
	pdfRenderer pdf.Renderer,
	excelRenderer *excel.Renderer,
) *Invoice {
	return &Invoice{
		generator:     generator,
		pdfRenderer:   pdfRenderer,
		excelRenderer: excelRenderer,
	}
}

func (h *Invoice) GenerateMonthly(c *gin.Context) {
	var req GenerateMonthlyInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	format := c.DefaultQuery("format", "pdf")

	cmd := service.GenerateMonthlyInvoiceCommand{
		Month:        req.Month,
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
	}

	// 1. Domain invoice
	invoice, err := h.generator.GenerateMonthly(
		c.Request.Context(),
		cmd,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch format {
	case "pdf":
		h.renderPDF(c, invoice)

	case "excel":
		h.renderExcel(c, invoice)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported format: " + format,
		})
	}
}

func (h *Invoice) renderPDF(
	c *gin.Context,
	invoice *domain.Invoice,
) {
	renderInvoice := render.BuildInvoice(*invoice)

	html, err := render.HTML(renderInvoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	footer, err := render.FooterHTML()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := h.pdfRenderer.RenderHTML(
		c.Request.Context(),
		string(html),
		string(footer),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := buildInvoiceFilename(renderInvoice.Consultant.FullName, renderInvoice.Invoice.Number, "pdf")

	c.Header(
		"Content-Disposition",
		`inline; filename="`+filename+`"`,
	)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Invoice) renderExcel(
	c *gin.Context,
	invoice *domain.Invoice,
) {
	fullName := consultantFullName(invoice.Consultant)
	filename := buildInvoiceFilename(fullName, invoice.Number, "xlsx")

	c.Header(
		"Content-Disposition",
		`attachment; filename="`+filename+`"`,
	)
	c.Header(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)

	excelBytes, err := h.excelRenderer.Render(*invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate excel invoice",
		})
		return
	}

	c.Data(
		http.StatusOK,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		excelBytes,
	)
}

func consultantFullName(c domain.Consultant) string {
	parts := []string{c.FirstName}

	if c.MiddleName != "" {
		parts = append(parts, c.MiddleName)
	}

	parts = append(parts, c.LastName)

	return strings.Join(parts, " ")
}

func buildInvoiceFilename(fullName string, invoiceNumber string, ext string) string {
	return fmt.Sprintf(
		"%s_Invoice_%s.%s",
		fullName,
		invoiceNumber,
		ext,
	)
}
