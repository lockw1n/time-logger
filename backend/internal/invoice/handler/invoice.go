package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lockw1n/time-logger/internal/invoice/render"
	"github.com/lockw1n/time-logger/internal/invoice/service"
)

type Invoice struct {
	generator service.InvoiceGenerator
}

func NewInvoice(generator service.InvoiceGenerator) *Invoice {
	return &Invoice{generator: generator}
}

func (h *Invoice) GenerateMonthly(c *gin.Context) {
	var req GenerateMonthlyInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := service.GenerateMonthlyInvoiceCommand{
		Month:        req.Month,
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
	}

	pdfBytes, renderInvoice, err := h.generator.GenerateMonthlyPDF(
		c.Request.Context(),
		cmd,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := buildInvoiceFilename(renderInvoice)

	c.Header(
		"Content-Disposition",
		`inline; filename="`+filename+`"`,
	)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func buildInvoiceFilename(invoice render.Invoice) string {
	return fmt.Sprintf(
		"%s_Invoice_%s.pdf",
		invoice.Consultant.FullName,
		invoice.Invoice.Number,
	)
}
