package repository

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/company/domain"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, company domain.Company) (domain.Company, error) {
	model := toModel(company)

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return domain.Company{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) Update(ctx context.Context, company domain.Company) (domain.Company, error) {
	model := toModel(company)

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return domain.Company{}, mapError(err)
	}

	var updated companyModel
	if err := r.db.WithContext(ctx).
		First(&updated, model.ID).
		Error; err != nil {
		return domain.Company{}, mapError(err)
	}

	return toDomain(updated), nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&companyModel{}, id)

	if err := result.Error; err != nil {
		return mapError(err)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (domain.Company, error) {
	var model companyModel

	if err := r.db.WithContext(ctx).
		First(&model, id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Company{}, ErrNotFound
		}

		return domain.Company{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) ListAll(ctx context.Context) ([]domain.Company, error) {
	var models []companyModel

	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}

func (r *gormRepository) ListByConsultant(ctx context.Context, consultantID uint64) ([]domain.Company, error) {
	var models []companyModel

	if err := r.db.WithContext(ctx).
		Joins("JOIN contracts ON contracts.company_id = companies.id").
		Where("contracts.consultant_id = ?", consultantID).
		Order("companies.name ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}
