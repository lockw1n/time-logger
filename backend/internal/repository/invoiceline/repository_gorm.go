package invoiceline

import (
	"github.com/lockw1n/time-logger/internal/models"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(entry *models.InvoiceLine) error {
	return r.db.Create(entry).Error
}

func (r *gormRepository) Delete(id uint64) error {
	return r.db.Delete(&models.InvoiceLine{}, id).Error
}

func (r *gormRepository) FindByInvoice(invoiceID uint64) ([]models.InvoiceLine, error) {
	var list []models.InvoiceLine
	err := r.db.
		Where("invoice_id = ?", invoiceID).
		Find(&list).Error
	return list, err
}

func (r *gormRepository) FindByEntry(entryID uint64) ([]models.InvoiceLine, error) {
	var list []models.InvoiceLine
	err := r.db.
		Where("entry_id = ?", entryID).
		Find(&list).Error
	return list, err
}
