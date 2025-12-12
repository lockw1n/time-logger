package label

import (
	"errors"

	"github.com/lockw1n/time-logger/internal/models"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(label *models.Label) (*models.Label, error) {
	if err := r.db.Create(label).Error; err != nil {
		return nil, err
	}

	var out models.Label
	if err := r.db.
		Preload("Company").
		First(&out, label.ID).
		Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *gormRepository) Update(label *models.Label) (*models.Label, error) {
	result := r.db.
		Model(&models.Label{}).
		Where("id = ?", label.ID).
		Updates(label)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	var updated models.Label
	if err := r.db.
		Preload("Company").
		First(&updated, label.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *gormRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Label{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(id uint64) (*models.Label, error) {
	var entry models.Label
	err := r.db.
		Preload("Company").
		First(&entry, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *gormRepository) FindByCompany(companyID uint64) ([]models.Label, error) {
	var list []models.Label
	err := r.db.
		Where("company_id = ?", companyID).
		Preload("Company").
		Find(&list).Error
	return list, err
}

func (r *gormRepository) FindByName(companyID uint64, name string) (*models.Label, error) {
	var label models.Label
	err := r.db.
		Where("company_id = ? AND LOWER(name) = LOWER(?)", companyID, name).
		Preload("Company").
		First(&label).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &label, nil
}
