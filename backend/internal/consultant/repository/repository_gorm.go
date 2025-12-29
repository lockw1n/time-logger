package repository

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/consultant/domain"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, consultant domain.Consultant) (domain.Consultant, error) {
	model := toModel(consultant)

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return domain.Consultant{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) Update(ctx context.Context, consultant domain.Consultant) (domain.Consultant, error) {
	model := toModel(consultant)

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return domain.Consultant{}, mapError(err)
	}

	var updated consultantModel
	if err := r.db.WithContext(ctx).
		First(&updated, model.ID).
		Error; err != nil {
		return domain.Consultant{}, mapError(err)
	}

	return toDomain(updated), nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&consultantModel{}, id)

	if err := result.Error; err != nil {
		return mapError(err)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (domain.Consultant, error) {
	var model consultantModel

	if err := r.db.WithContext(ctx).
		First(&model, id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Consultant{}, ErrNotFound
		}

		return domain.Consultant{}, mapError(err)
	}

	return toDomain(model), nil
}
