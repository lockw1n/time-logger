package entry

import (
	"errors"
	"time"

	"github.com/lockw1n/time-logger/internal/models"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(entry *models.Entry) (*models.Entry, error) {
	if err := r.db.Create(entry).Error; err != nil {
		return nil, err
	}

	var out models.Entry
	if err := r.db.
		Preload("Ticket").
		Preload("Label").
		First(&out, entry.ID).
		Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *gormRepository) Update(entry *models.Entry) (*models.Entry, error) {
	result := r.db.
		Model(&models.Entry{}).
		Where("id = ?", entry.ID).
		Updates(entry)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	var updated models.Entry
	if err := r.db.
		Preload("Ticket").
		Preload("Label").
		First(&updated, entry.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *gormRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Entry{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(id uint64) (*models.Entry, error) {
	var entry models.Entry
	err := r.db.
		Preload("Ticket").
		Preload("Label").
		Preload("Consultant").
		Preload("Company").
		First(&entry, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *gormRepository) FindByCompany(companyID uint64) ([]models.Entry, error) {
	var entries []models.Entry
	err := r.db.
		Where("company_id = ?", companyID).
		Preload("Ticket").
		Preload("Label").
		Preload("Consultant").
		Preload("Company").
		Find(&entries).Error

	return entries, err
}

func (r *gormRepository) FindByConsultant(consultantID uint64) ([]models.Entry, error) {
	var entries []models.Entry
	err := r.db.
		Where("consultant_id = ?", consultantID).
		Preload("Ticket").
		Preload("Label").
		Preload("Consultant").
		Preload("Company").
		Find(&entries).Error

	return entries, err
}

func (r *gormRepository) ListAll() ([]models.Entry, error) {
	var entries []models.Entry
	err := r.db.
		Preload("Ticket").
		Preload("Label").
		Preload("Consultant").
		Preload("Company").
		Find(&entries).Error

	return entries, err
}

func (r *gormRepository) FindWithDetails(
	consultantID uint64,
	companyID uint64,
	start time.Time,
	end time.Time,
) ([]models.Entry, error) {
	var entries []models.Entry

	err := r.db.
		Preload("Ticket").
		Preload("Label").
		Where("consultant_id = ?", consultantID).
		Where("company_id = ?", companyID).
		Where("date BETWEEN ? AND ?", start, end).
		Order("date ASC").
		Find(&entries).Error

	return entries, err
}
