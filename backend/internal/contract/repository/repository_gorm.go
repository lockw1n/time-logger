package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lockw1n/time-logger/internal/contract/domain"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, contract domain.Contract) (domain.Contract, error) {
	model := toModel(contract)

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return domain.Contract{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) Update(ctx context.Context, contract domain.Contract) (domain.Contract, error) {
	model := toModel(contract)

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return domain.Contract{}, mapError(err)
	}

	var updated contractModel
	if err := r.db.WithContext(ctx).
		First(&updated, model.ID).
		Error; err != nil {
		return domain.Contract{}, mapError(err)
	}

	return toDomain(updated), nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&contractModel{}, id)

	if err := result.Error; err != nil {
		return mapError(err)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (domain.Contract, error) {
	var model contractModel

	if err := r.db.WithContext(ctx).
		First(&model, id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Contract{}, ErrNotFound
		}

		return domain.Contract{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) ListByConsultant(ctx context.Context, consultantID uint64) ([]domain.Contract, error) {
	var models []contractModel

	if err := r.db.WithContext(ctx).
		Where("consultant_id = ?", consultantID).
		Order("created_at ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}

func (r *gormRepository) FindActiveByConsultantCompany(
	ctx context.Context,
	consultantID uint64,
	companyID uint64,
	at time.Time,
) (domain.Contract, error) {
	var model contractModel

	if err := r.db.WithContext(ctx).
		Where("consultant_id = ?", consultantID).
		Where("company_id = ?", companyID).
		Where("start_date <= ?", at).
		Where("(end_date IS NULL OR end_date >= ?)", at).
		First(&model).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Contract{}, ErrNotFound
		}

		return domain.Contract{}, mapError(err)
	}

	return toDomain(model), nil
}
