package company

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

func (r *gormRepository) Create(company *models.Company) (*models.Company, error) {
	if err := r.db.Create(company).Error; err != nil {
		return nil, err
	}

	var out models.Company
	if err := r.db.
		First(&out, company.ID).
		Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *gormRepository) Update(company *models.Company) (*models.Company, error) {
	result := r.db.
		Model(&models.Company{}).
		Where("id = ?", company.ID).
		Updates(company)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	var updated models.Company
	if err := r.db.First(&updated, company.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *gormRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Company{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(id uint64) (*models.Company, error) {
	var company models.Company
	err := r.db.First(&company, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *gormRepository) FindByName(name string) (*models.Company, error) {
	var company models.Company
	err := r.db.
		Where("LOWER(name) = LOWER(?)", name).
		First(&company).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *gormRepository) ListAll() ([]models.Company, error) {
	var companies []models.Company
	err := r.db.
		Order("name ASC").
		Find(&companies).Error

	return companies, err
}
