package repository

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/activity/domain"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	model := toModel(activity)

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return domain.Activity{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) Update(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	model := toModel(activity)

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return domain.Activity{}, mapError(err)
	}

	var updated activityModel
	if err := r.db.WithContext(ctx).
		First(&updated, model.ID).
		Error; err != nil {
		return domain.Activity{}, mapError(err)
	}

	return toDomain(updated), nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&activityModel{}, id)

	if err := result.Error; err != nil {
		return mapError(err)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (domain.Activity, error) {
	var model activityModel

	if err := r.db.WithContext(ctx).
		First(&model, id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Activity{}, ErrNotFound
		}

		return domain.Activity{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) ListByIDs(ctx context.Context, ids []uint64) ([]domain.Activity, error) {
	if len(ids) == 0 {
		return []domain.Activity{}, nil
	}

	var models []activityModel
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Order("priority ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}

func (r *gormRepository) ListByCompany(ctx context.Context, companyID uint64) ([]domain.Activity, error) {
	var models []activityModel

	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("priority ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}
