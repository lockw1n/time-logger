package repository

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/ticket/domain"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error) {
	model := toModel(ticket)

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return domain.Ticket{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) Update(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error) {
	model := toModel(ticket)

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return domain.Ticket{}, mapError(err)
	}

	var updated ticketModel
	if err := r.db.WithContext(ctx).
		First(&updated, model.ID).
		Error; err != nil {
		return domain.Ticket{}, mapError(err)
	}

	return toDomain(updated), nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&ticketModel{}, id)

	if err := result.Error; err != nil {
		return mapError(err)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uint64) (domain.Ticket, error) {
	var model ticketModel

	if err := r.db.WithContext(ctx).
		First(&model, id).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Ticket{}, ErrNotFound
		}

		return domain.Ticket{}, mapError(err)
	}

	return toDomain(model), nil
}

func (r *gormRepository) ListByIDs(ctx context.Context, ids []uint64) ([]domain.Ticket, error) {
	if len(ids) == 0 {
		return []domain.Ticket{}, nil
	}

	var models []ticketModel
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Order("code ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}

func (r *gormRepository) ListByCompany(ctx context.Context, companyID uint64) ([]domain.Ticket, error) {
	var models []ticketModel

	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("code ASC").
		Find(&models).
		Error; err != nil {
		return nil, mapError(err)
	}

	return toDomainSlice(models), nil
}

func (r *gormRepository) FindByCompanyAndCode(ctx context.Context, companyID uint64, code string) (domain.Ticket, error) {
	var model ticketModel

	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Where("code = ?", code).
		First(&model).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Ticket{}, ErrNotFound
		}

		return domain.Ticket{}, mapError(err)
	}

	return toDomain(model), nil
}
