package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/invoice/domain"
)

type Service interface {
	GenerateInvoice(ctx context.Context, input GenerateInvoiceInput) (domain.Invoice, error)
}
