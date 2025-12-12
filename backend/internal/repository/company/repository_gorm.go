package company

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

func (r *gormRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *gormRepository) Update(company *models.Company) error {
	return r.db.Save(company).Error
}

func (r *gormRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Company{}, id).Error
}

func (r *gormRepository) FindByID(id uint64) (*models.Company, error) {
	var company models.Company
	err := r.db.First(&company, id).Error
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
