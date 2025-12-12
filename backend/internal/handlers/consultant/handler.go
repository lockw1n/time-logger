package consultant

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	consultantdto "github.com/lockw1n/time-logger/internal/dto/consultant"
	consultantservice "github.com/lockw1n/time-logger/internal/service/consultant"
)

type Handler struct {
	service consultantservice.Service
}

func NewHandler(service consultantservice.Service) *Handler {
	return &Handler{service: service}
}

var ErrInvalidConsultant = errors.New("invalid consultant data")

/************************************************
 *               Helpers
 ************************************************/

func normalizeInput(in consultantdto.Request) consultantdto.Request {
	return consultantdto.Request{
		FirstName:    strings.TrimSpace(in.FirstName),
		MiddleName:   trimPtr(in.MiddleName),
		LastName:     strings.TrimSpace(in.LastName),
		AddressLine1: strings.TrimSpace(in.AddressLine1),
		AddressLine2: trimPtr(in.AddressLine2),
		Zip:          strings.TrimSpace(in.Zip),
		City:         strings.TrimSpace(in.City),
		Region:       trimPtr(in.Region),
		Country:      strings.TrimSpace(in.Country),
		TaxNumber:    trimPtr(in.TaxNumber),
		BankName:     strings.TrimSpace(in.BankName),
		BankAddress:  strings.TrimSpace(in.BankAddress),
		BankCountry:  strings.TrimSpace(in.BankCountry),
		BankIBAN:     strings.TrimSpace(in.BankIBAN),
		BankBIC:      strings.TrimSpace(in.BankBIC),
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

func validateInput(in consultantdto.Request) error {
	if in.FirstName == "" ||
		in.LastName == "" ||
		in.Zip == "" ||
		in.Country == "" ||
		in.City == "" ||
		in.AddressLine1 == "" ||
		in.BankName == "" ||
		in.BankAddress == "" ||
		in.BankCountry == "" ||
		in.BankIBAN == "" ||
		in.BankBIC == "" {
		return ErrInvalidConsultant
	}
	return nil
}

/************************************************
 *                  Routes
 ************************************************/

// GetConsultant GET /consultant
func (h *Handler) GetConsultant(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	consultant, err := h.service.Get(id)
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

// UpsertConsultant PUT /consultant   (same API as before)
func (h *Handler) UpsertConsultant(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req consultantdto.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req = normalizeInput(req)

	if err := validateInput(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "first_name, last_name, country, city, and address_line1 are required",
		})
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save consultant"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
