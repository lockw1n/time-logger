package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/lockw1n/time-logger/internal/models"
)

var (
	ErrInvalidCompany = errors.New("invalid company data")
)

type companyInput struct {
	Name         string `json:"name"`
	UID          string `json:"uid"`
	AddressLine1 string `json:"address_line1"`
	Zip          string `json:"zip"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Payment      string `json:"payment_condition"`
}

type CompanyService struct {
	DB *gorm.DB
}

func NewCompanyService(db *gorm.DB) *CompanyService {
	return &CompanyService{DB: db}
}

func normalizeCompanyInput(in companyInput) companyInput {
	return companyInput{
		Name:         strings.TrimSpace(in.Name),
		UID:          strings.TrimSpace(in.UID),
		AddressLine1: strings.TrimSpace(in.AddressLine1),
		Zip:          strings.TrimSpace(in.Zip),
		City:         strings.TrimSpace(in.City),
		Country:      strings.TrimSpace(in.Country),
		Payment:      strings.TrimSpace(in.Payment),
	}
}

func validateCompanyInput(in companyInput) error {
	if in.Name == "" || in.AddressLine1 == "" || in.Zip == "" || in.City == "" || in.Country == "" {
		return ErrInvalidCompany
	}
	return nil
}

// Get returns the single company record, or nil if none exists.
func (s *CompanyService) Get() (*models.Company, error) {
	var company models.Company
	err := s.DB.First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

// Upsert creates or updates the single company record.
func (s *CompanyService) Upsert(in companyInput) (*models.Company, error) {
	in = normalizeCompanyInput(in)
	if err := validateCompanyInput(in); err != nil {
		return nil, err
	}

	var company models.Company
	err := s.DB.First(&company).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		company = models.Company{
			Name:         in.Name,
			UID:          in.UID,
			AddressLine1: in.AddressLine1,
			Zip:          in.Zip,
			City:         in.City,
			Country:      in.Country,
			Payment:      in.Payment,
		}
		if err := s.DB.Create(&company).Error; err != nil {
			return nil, err
		}
		return &company, nil
	}
	if err != nil {
		return nil, err
	}

	company.Name = in.Name
	company.UID = in.UID
	company.AddressLine1 = in.AddressLine1
	company.Zip = in.Zip
	company.City = in.City
	company.Country = in.Country
	company.Payment = in.Payment

	if err := s.DB.Save(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

type CompanyHandler struct {
	Service *CompanyService
}

func NewCompanyHandler(service *CompanyService) *CompanyHandler {
	return &CompanyHandler{Service: service}
}

// GetCompany returns the stored company details.
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	company, err := h.Service.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch company"})
		return
	}
	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not configured"})
		return
	}
	c.JSON(http.StatusOK, company)
}

// UpsertCompany creates or updates the company record.
func (h *CompanyHandler) UpsertCompany(c *gin.Context) {
	var input companyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	company, err := h.Service.Upsert(input)
	if err != nil {
		if errors.Is(err, ErrInvalidCompany) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, address_line1, zip, city, and country are required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save company"})
		return
	}

	c.JSON(http.StatusOK, company)
}
