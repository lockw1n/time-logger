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
	ErrInvalidConsultant = errors.New("invalid consultant data")
)

type consultantInput struct {
	FirstName   string  `json:"first_name"`
	MiddleName  string  `json:"middle_name"`
	LastName    string  `json:"last_name"`
	Country     string  `json:"country"`
	Zip         string  `json:"zip"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Address1    string  `json:"address_line1"`
	Address2    string  `json:"address_line2"`
	TaxNumber   string  `json:"tax_number"`
	BankName    string  `json:"bank_name"`
	BankAddress string  `json:"bank_address"`
	BankCountry string  `json:"bank_country"`
	IBAN        string  `json:"iban"`
	BIC         string  `json:"bic"`
	OrderNumber string  `json:"order_number"`
	HourlyRate  float64 `json:"hourly_rate"`
}

type ConsultantService struct {
	DB *gorm.DB
}

func NewConsultantService(db *gorm.DB) *ConsultantService {
	return &ConsultantService{DB: db}
}

func normalizeConsultantInput(in consultantInput) consultantInput {
	return consultantInput{
		FirstName:   strings.TrimSpace(in.FirstName),
		MiddleName:  strings.TrimSpace(in.MiddleName),
		LastName:    strings.TrimSpace(in.LastName),
		Country:     strings.TrimSpace(in.Country),
		Zip:         strings.TrimSpace(in.Zip),
		Region:      strings.TrimSpace(in.Region),
		City:        strings.TrimSpace(in.City),
		Address1:    strings.TrimSpace(in.Address1),
		Address2:    strings.TrimSpace(in.Address2),
		TaxNumber:   strings.TrimSpace(in.TaxNumber),
		BankName:    strings.TrimSpace(in.BankName),
		BankAddress: strings.TrimSpace(in.BankAddress),
		BankCountry: strings.TrimSpace(in.BankCountry),
		IBAN:        strings.TrimSpace(in.IBAN),
		BIC:         strings.TrimSpace(in.BIC),
		OrderNumber: strings.TrimSpace(in.OrderNumber),
		HourlyRate:  in.HourlyRate,
	}
}

func validateConsultantInput(in consultantInput) error {
	if in.FirstName == "" || in.LastName == "" || in.Country == "" || in.City == "" || in.Address1 == "" {
		return ErrInvalidConsultant
	}
	return nil
}

func (s *ConsultantService) Get() (*models.Consultant, error) {
	var consultant models.Consultant
	err := s.DB.First(&consultant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &consultant, nil
}

func (s *ConsultantService) Upsert(in consultantInput) (*models.Consultant, error) {
	in = normalizeConsultantInput(in)
	if err := validateConsultantInput(in); err != nil {
		return nil, err
	}

	var consultant models.Consultant
	err := s.DB.First(&consultant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		consultant = models.Consultant{
			FirstName:   in.FirstName,
			MiddleName:  in.MiddleName,
			LastName:    in.LastName,
			Country:     in.Country,
			Zip:         in.Zip,
			Region:      in.Region,
			City:        in.City,
			Address1:    in.Address1,
			Address2:    in.Address2,
			TaxNumber:   in.TaxNumber,
			BankName:    in.BankName,
			BankAddress: in.BankAddress,
			BankCountry: in.BankCountry,
			IBAN:        in.IBAN,
			BIC:         in.BIC,
			OrderNumber: in.OrderNumber,
			HourlyRate:  in.HourlyRate,
		}
		if err := s.DB.Create(&consultant).Error; err != nil {
			return nil, err
		}
		return &consultant, nil
	}
	if err != nil {
		return nil, err
	}

	consultant.FirstName = in.FirstName
	consultant.MiddleName = in.MiddleName
	consultant.LastName = in.LastName
	consultant.Country = in.Country
	consultant.Zip = in.Zip
	consultant.Region = in.Region
	consultant.City = in.City
	consultant.Address1 = in.Address1
	consultant.Address2 = in.Address2
	consultant.TaxNumber = in.TaxNumber
	consultant.BankName = in.BankName
	consultant.BankAddress = in.BankAddress
	consultant.BankCountry = in.BankCountry
	consultant.IBAN = in.IBAN
	consultant.BIC = in.BIC
	consultant.OrderNumber = in.OrderNumber
	consultant.HourlyRate = in.HourlyRate

	if err := s.DB.Save(&consultant).Error; err != nil {
		return nil, err
	}
	return &consultant, nil
}

type ConsultantHandler struct {
	Service *ConsultantService
}

func NewConsultantHandler(service *ConsultantService) *ConsultantHandler {
	return &ConsultantHandler{Service: service}
}

func (h *ConsultantHandler) GetConsultant(c *gin.Context) {
	consultant, err := h.Service.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch consultant"})
		return
	}
	if consultant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "consultant not configured"})
		return
	}
	c.JSON(http.StatusOK, consultant)
}

func (h *ConsultantHandler) UpsertConsultant(c *gin.Context) {
	var input consultantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	consultant, err := h.Service.Upsert(input)
	if err != nil {
		if err == ErrInvalidConsultant {
			c.JSON(http.StatusBadRequest, gin.H{"error": "first_name, last_name, country, city, and address_line1 are required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save consultant"})
		return
	}

	c.JSON(http.StatusOK, consultant)
}
