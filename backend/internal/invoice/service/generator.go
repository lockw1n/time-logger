package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/invoice/domain"
)

type InvoiceGenerator interface {
	GenerateMonthly(ctx context.Context, cmd GenerateMonthlyInvoiceCommand) (*domain.Invoice, error)
}
