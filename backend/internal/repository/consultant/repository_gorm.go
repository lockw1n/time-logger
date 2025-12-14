package consultant

import (
	"context"
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

func (r *gormRepository) Create(consultant *models.Consultant) (*models.Consultant, error) {
	if err := r.db.Create(consultant).Error; err != nil {
		return nil, err
	}

	var out models.Consultant
	if err := r.db.
		First(&out, consultant.ID).
		Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *gormRepository) Update(consultant *models.Consultant) (*models.Consultant, error) {
	result := r.db.
		Model(&models.Consultant{}).
		Where("id = ?", consultant.ID).
		Updates(consultant)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	var updated models.Consultant
	if err := r.db.First(&updated, consultant.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *gormRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Consultant{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (*models.Consultant, error) {
	var cons models.Consultant
	err := r.db.
		WithContext(ctx).
		First(&cons, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

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
