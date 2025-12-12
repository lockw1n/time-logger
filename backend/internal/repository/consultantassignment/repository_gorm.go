package consultantassignment

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

func (r *gormRepository) Create(assignment *models.ConsultantAssignment) (*models.ConsultantAssignment, error) {
	if err := r.db.Create(assignment).Error; err != nil {
		return nil, err
	}

	// Reload with preloads for response mapping
	var out models.ConsultantAssignment
	if err := r.db.
		Preload("Ticket").
		Preload("Label").
		First(&out, assignment.ID).
		Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *gormRepository) Update(assignment *models.ConsultantAssignment) (*models.ConsultantAssignment, error) {
	result := r.db.
		Model(&models.ConsultantAssignment{}).
		Where("id = ?", assignment.ID).
		Updates(assignment)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Reload entry with preloads
	var updated models.ConsultantAssignment
	if err := r.db.
		Preload("Consultant").
		Preload("Company").
		First(&updated, assignment.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *gormRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.ConsultantAssignment{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) FindByID(id uint64) (*models.ConsultantAssignment, error) {
	var cc models.ConsultantAssignment
	err := r.db.
		Preload("Consultant").
		Preload("Company").
		First(&cc, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &cc, nil
}

func (r *gormRepository) FindByPair(consultantID, companyID uint64) (*models.ConsultantAssignment, error) {
	var cc models.ConsultantAssignment
	err := r.db.
		Where("consultant_id = ? AND company_id = ?", consultantID, companyID).
		Preload("Consultant").
		Preload("Company").
		First(&cc).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return &cc, nil
}

func (r *gormRepository) FindByConsultant(consultantID uint64) ([]models.ConsultantAssignment, error) {
	var list []models.ConsultantAssignment
	err := r.db.
		Where("consultant_id = ?", consultantID).
		Preload("Consultant").
		Preload("Company").
		Find(&list).Error

	return list, err
}

func (r *gormRepository) FindByCompany(companyID uint64) ([]models.ConsultantAssignment, error) {
	var list []models.ConsultantAssignment
	err := r.db.
		Where("company_id = ?", companyID).
		Preload("Consultant").
		Preload("Company").
		Find(&list).Error

	return list, err
}
