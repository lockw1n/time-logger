package company

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	companydto "github.com/lockw1n/time-logger/internal/dto/company"
	companyservice "github.com/lockw1n/time-logger/internal/service/company"
)

type Handler struct {
	service companyservice.Service
}

func NewHandler(service companyservice.Service) *Handler {
	return &Handler{service: service}
}

var ErrInvalidCompany = errors.New("invalid company data")

/******************** Helpers ********************/

func normalizeInput(in companydto.Request) companydto.Request {
	return companydto.Request{
		Name:         strings.TrimSpace(in.Name),
		NameShort:    trimPtr(in.NameShort),
		TaxNumber:    trimPtr(in.TaxNumber),
		AddressLine1: strings.TrimSpace(in.AddressLine1),
		AddressLine2: trimPtr(in.AddressLine2),
		Zip:          strings.TrimSpace(in.Zip),
		City:         strings.TrimSpace(in.City),
		Region:       trimPtr(in.Region),
		Country:      strings.TrimSpace(in.Country),
	}
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

func parseIDParam(c *gin.Context) (uint64, bool) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func validateInput(in companydto.Request) error {
	if in.Name == "" ||
		in.TaxNumber == nil || *in.TaxNumber == "" ||
		in.AddressLine1 == "" ||
		in.Zip == "" ||
		in.City == "" ||
		in.Country == "" {
		return ErrInvalidCompany
	}

	return nil
}

/******************** Routes ********************/

// GetCompany GET /company
func (h *Handler) GetCompany(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	company, err := h.service.Get(c.Request.Context(), id)
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

// UpsertCompany PUT /company  (or POST, depending on your router)
func (h *Handler) UpsertCompany(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req companydto.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req = normalizeInput(req)

	if err := validateInput(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name, tax_number, address_line1, zip, city, and country are required",
		})
		return
	}

	resp, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save company"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
