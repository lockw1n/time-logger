package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lockw1n/time-logger/internal/entry/domain"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, entry domain.Entry) (domain.Entry, error) {
	model := toModel(entry)

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return domain.Entry{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) Update(ctx context.Context, entry domain.Entry) (domain.Entry, error) {
	model := toModel(entry)

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return domain.Entry{}, mapError(err)
	}

	var updated entryModel
	if err := r.db.WithContext(ctx).
		First(&updated, model.ID).
		Error; err != nil {
		return domain.Entry{}, mapError(err)
	}

	return toDomain(updated), nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&entryModel{}, id)

	if err := result.Error; err != nil {
		return mapError(err)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (domain.Entry, error) {
	var model entryModel

	if err := r.db.WithContext(ctx).
		First(&model, id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Entry{}, ErrNotFound
		}

		return domain.Entry{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) ListForConsultantPeriod(
	ctx context.Context,
	consultantID uint64,
	companyID uint64,
	start time.Time,
	end time.Time,
) ([]domain.Entry, error) {
	var models []entryModel

	if err := r.db.WithContext(ctx).
		Where("consultant_id = ?", consultantID).
		Where("company_id = ?", companyID).
		Where("date BETWEEN ? AND ?", start, end).
		Order("date ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}
