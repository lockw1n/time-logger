package invoiceline

import (
	"time"

	invoicelinedto "github.com/lockw1n/time-logger/internal/dto/invoiceline"
	"github.com/lockw1n/time-logger/internal/models"
)

/* ============================================================
   MODEL → DTO (Single Invoice Line)
=============================================================== */

func ToLine(line *models.InvoiceLine) invoicelinedto.Line {
	return invoicelinedto.Line{
		ID:        line.ID,
		EntryID:   line.EntryID,
		Hours:     line.Hours,
		Rate:      line.Rate,
		Amount:    line.Amount,
		CreatedAt: line.CreatedAt.Format(time.RFC3339),
	}
}

/* ============================================================
   MODEL SLICE → DTO SLICE (Batch)
=============================================================== */

func ToLines(lines []models.InvoiceLine) []invoicelinedto.Line {
	out := make([]invoicelinedto.Line, len(lines))
	for i := range lines {
		out[i] = ToLine(&lines[i])
	}
	return out
}
