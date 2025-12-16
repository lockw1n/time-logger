package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/invoice/render"
)

type InvoiceGenerator interface {
	GenerateMonthlyPDF(ctx context.Context, cmd GenerateMonthlyInvoiceCommand) ([]byte, render.Invoice, error)
}
