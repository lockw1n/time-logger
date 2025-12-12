package invoice

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

func (r *gormRepository) Create(inv *models.Invoice) error {
	return r.db.Create(inv).Error
}

func (r *gormRepository) Update(inv *models.Invoice) error {
	return r.db.Save(inv).Error
}

func (r *gormRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Invoice{}, id).Error
}

func (r *gormRepository) FindByID(id uint64) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.
		Preload("Entries").
		First(&inv, id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *gormRepository) FindByCompany(companyID uint64) ([]models.Invoice, error) {
	var list []models.Invoice
	err := r.db.
		Where("company_id = ?", companyID).
		Find(&list).Error
	return list, err
}

func (r *gormRepository) FindByConsultant(consultantID uint64) ([]models.Invoice, error) {
	var list []models.Invoice
	err := r.db.
		Where("consultant_id = ?", consultantID).
		Find(&list).Error
	return list, err
}

func (r *gormRepository) ListAll() ([]models.Invoice, error) {
	var list []models.Invoice
	err := r.db.Find(&list).Error
	return list, err
}
