package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/company/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeCreateCompanyRequest(&req)

	input := service.CreateCompanyInput{
		Name:         req.Name,
		NameShort:    req.NameShort,
		TaxNumber:    req.TaxNumber,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Zip:          req.Zip,
		City:         req.City,
		Region:       req.Region,
		Country:      req.Country,
	}

	company, err := h.service.CreateCompany(c.Request.Context(), input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, toResponse(company))
}

func (h *Handler) UpdateCompany(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeUpdateCompanyRequest(&req)

	input := service.UpdateCompanyInput{
		Name:         req.Name,
		NameShort:    req.NameShort,
		TaxNumber:    req.TaxNumber,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Zip:          req.Zip,
		City:         req.City,
		Region:       req.Region,
		Country:      req.Country,
	}

	company, err := h.service.UpdateCompany(c.Request.Context(), id, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(company))
}

func (h *Handler) DeleteCompany(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteCompany(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetCompany(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	company, err := h.service.GetCompany(c.Request.Context(), id)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(company))
}

func (h *Handler) ListCompanies(c *gin.Context) {
	ctx := c.Request.Context()

	if consultantIDStr := c.Query("consultant_id"); consultantIDStr != "" {
		consultantID, err := strconv.ParseUint(consultantIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid consultant_id"})
			return
		}

		companies, err := h.service.ListCompaniesForConsultant(ctx, consultantID)
		if err != nil {
			status, msg := mapError(err)
			c.JSON(status, gin.H{"error": msg})
			return
		}

		respondCompanies(c, companies)
		return
	}

	companies, err := h.service.ListCompanies(ctx)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	respondCompanies(c, companies)
}
