package label

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

func (r *gormRepository) Create(label *models.Label) error {
	return r.db.Create(label).Error
}

func (r *gormRepository) Update(label *models.Label) error {
	return r.db.Save(label).Error
}

func (r *gormRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Label{}, id).Error
}

func (r *gormRepository) FindByID(id uint64) (*models.Label, error) {
	var label models.Label
	err := r.db.Preload("Company").First(&label, id).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *gormRepository) FindByCompany(companyID uint64) ([]models.Label, error) {
	var list []models.Label
	err := r.db.Where("company_id = ?", companyID).Find(&list).Error
	return list, err
}

func (r *gormRepository) FindByName(companyID uint64, name string) (*models.Label, error) {
	var label models.Label
	err := r.db.
		Where("company_id = ? AND LOWER(name) = LOWER(?)", companyID, name).
		First(&label).Error

	if err != nil {
		return nil, err
	}
	return &label, nil
}
