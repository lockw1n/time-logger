package invoiceline

import "github.com/lockw1n/time-logger/internal/models"

type Repository interface {
	Create(entry *models.InvoiceLine) error
	Delete(id uint64) error

	FindByInvoice(invoiceID uint64) ([]models.InvoiceLine, error)
	FindByEntry(entryID uint64) ([]models.InvoiceLine, error)
}
