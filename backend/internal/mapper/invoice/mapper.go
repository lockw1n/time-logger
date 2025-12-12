package invoice

import (
	"time"

	invoicedto "github.com/lockw1n/time-logger/internal/dto/invoice"
	invoicelinedto "github.com/lockw1n/time-logger/internal/dto/invoiceline"
	"github.com/lockw1n/time-logger/internal/mapper/invoiceline"
	"github.com/lockw1n/time-logger/internal/models"
)

/* ============================================================
   STORED INVOICE (DETAILED) → DTO
=============================================================== */

func ToResponse(inv *models.Invoice) invoicedto.Response {
	var total float64
	if inv.TotalAmount != nil {
		total = *inv.TotalAmount
	}

	var metadata map[string]any
	if inv.Metadata != nil {
		metadata = inv.Metadata
	}

	return invoicedto.Response{
		ID:            inv.ID,
		ConsultantID:  inv.ConsultantID,
		CompanyID:     inv.CompanyID,
		PeriodStart:   inv.PeriodStart.Format("2006-01-02"),
		PeriodEnd:     inv.PeriodEnd.Format("2006-01-02"),
		TotalAmount:   total,
		Currency:      inv.Currency,
		Status:        inv.Status,
		InvoiceNumber: inv.InvoiceNumber,
		OrderNumber:   inv.OrderNumber,
		Metadata:      metadata,
		Lines:         invoiceline.ToLines(inv.Entries),
		CreatedAt:     inv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     inv.UpdatedAt.Format(time.RFC3339),
	}
}

/* ============================================================
   INVOICE LIST (SUMMARY VIEW) → DTO
=============================================================== */

func ToListItems(list []models.Invoice) []invoicedto.ListItem {
	out := make([]invoicedto.ListItem, len(list))

	for i := range list {
		inv := &list[i]

		var total float64
		if inv.TotalAmount != nil {
			total = *inv.TotalAmount
		}

		out[i] = invoicedto.ListItem{
			ID:            inv.ID,
			ConsultantID:  inv.ConsultantID,
			CompanyID:     inv.CompanyID,
			PeriodStart:   inv.PeriodStart.Format("2006-01-02"),
			PeriodEnd:     inv.PeriodEnd.Format("2006-01-02"),
			TotalAmount:   total,
			Currency:      inv.Currency,
			Status:        inv.Status,
			InvoiceNumber: inv.InvoiceNumber,
			OrderNumber:   inv.OrderNumber,
			CreatedAt:     inv.CreatedAt.Format(time.RFC3339),
		}
	}

	return out
}

/* ============================================================
   DYNAMICALLY GENERATED INVOICE (PREVIEW)
=============================================================== */

func ToGenerateResponse(
	consultantID uint64,
	companyID uint64,
	periodStart, periodEnd string,
	total float64,
	currency string,
	lines []invoicelinedto.Line,
) invoicedto.GenerateResponse {

	return invoicedto.GenerateResponse{
		ConsultantID: consultantID,
		CompanyID:    companyID,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		TotalAmount:  total,
		Currency:     currency,
		Lines:        lines,
	}
}
