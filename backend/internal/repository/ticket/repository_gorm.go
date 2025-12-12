package ticket

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

func (r *gormRepository) Create(ticket *models.Ticket) (*models.Ticket, error) {
	if err := r.db.Create(ticket).Error; err != nil {
		return nil, err
	}

	// Reload with preloads for response mapping
	var out models.Ticket
	if err := r.db.
		Preload("Company").
		First(&out, ticket.ID).
		Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *gormRepository) Update(ticket *models.Ticket) (*models.Ticket, error) {
	result := r.db.
		Model(&models.Ticket{}).
		Where("id = ?", ticket.ID).
		Updates(ticket)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Reload entry with preloads
	var updated models.Ticket
	if err := r.db.
		Preload("Company").
		First(&updated, ticket.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *gormRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Ticket{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(id uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.
		Preload("Company").
		First(&ticket, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *gormRepository) FindByCode(companyID uint64, code string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.
		Where("company_id = ? AND code = ?", companyID, code).
		First(&ticket).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *gormRepository) FindByCompany(companyID uint64) ([]models.Ticket, error) {
	var list []models.Ticket
	err := r.db.
		Where("company_id = ?", companyID).
		Find(&list).Error
	return list, err
}

func (r *gormRepository) ListAll() ([]models.Ticket, error) {
	var list []models.Ticket
	err := r.db.Find(&list).Error
	return list, err
}
