package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/invoice/domain"
	"github.com/lockw1n/time-logger/internal/invoice/render"
	excelrenderer "github.com/lockw1n/time-logger/internal/invoice/render/excel"
	htmlrenderer "github.com/lockw1n/time-logger/internal/invoice/render/html"
	"github.com/lockw1n/time-logger/internal/invoice/service"
	"github.com/lockw1n/time-logger/internal/pdf"
)

type Handler struct {
	service     service.Service
	pdfRenderer pdf.Renderer
}

func NewHandler(
	service service.Service,
	pdfRenderer pdf.Renderer,
) *Handler {
	return &Handler{
		service:     service,
		pdfRenderer: pdfRenderer,
	}
}

func (h *Handler) GenerateInvoice(c *gin.Context) {
	var req GenerateInvoiceRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	start, err := parseDate(req.Start)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	end, err := parseDate(req.End)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	format := strings.ToLower(c.DefaultQuery("format", "pdf"))

	input := service.GenerateInvoiceInput{
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
		Start:        start,
		End:          end,
	}

	invoice, err := h.service.GenerateInvoice(c.Request.Context(), input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
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

func (h *Handler) renderPDF(c *gin.Context, invoice domain.Invoice) {
	view := render.MapInvoiceToHTMLView(invoice)

	html, err := htmlrenderer.HTML(view)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rendering failed"})
		return
	}

	footer, err := htmlrenderer.FooterHTML()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := h.pdfRenderer.RenderHTML(c.Request.Context(), string(html), string(footer))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := buildInvoiceFilename(view.Consultant.FullName, view.Number, "pdf")

	c.Header(
		"Content-Disposition",
		buildContentDisposition(DispositionInline, filename),
	)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) renderExcel(c *gin.Context, invoice domain.Invoice) {
	view := render.MapInvoiceToExcelView(invoice)
	fullName := render.ConsultantFullName(invoice.Consultant)
	filename := buildInvoiceFilename(fullName, view.Number, "xlsx")

	c.Header(
		"Content-Disposition",
		buildContentDisposition(DispositionAttachment, filename),
	)

	excelBytes, err := excelrenderer.Render(view)
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
