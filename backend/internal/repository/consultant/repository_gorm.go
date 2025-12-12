package consultant

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

func (r *gormRepository) Create(cons *models.Consultant) error {
	return r.db.Create(cons).Error
}

func (r *gormRepository) Update(cons *models.Consultant) error {
	return r.db.Save(cons).Error
}

func (r *gormRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Consultant{}, id).Error
}

func (r *gormRepository) FindByID(id uint64) (*models.Consultant, error) {
	var cons models.Consultant
	err := r.db.First(&cons, id).Error
	if err != nil {
		return nil, err
	}
	return &cons, nil
}

func (r *gormRepository) FindByLastName(lastName string) ([]models.Consultant, error) {
	var result []models.Consultant
	err := r.db.
		Where("last_name ILIKE ?", "%"+lastName+"%").
		Find(&result).Error

	return result, err
}

func (r *gormRepository) ListAll() ([]models.Consultant, error) {
	var result []models.Consultant
	err := r.db.
		Order("last_name ASC, first_name ASC").
		Find(&result).Error

	return result, err
}
