package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/consultant/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateConsultant(c *gin.Context) {
	var req CreateConsultantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeCreateConsultantRequest(&req)

	input := service.CreateConsultantInput{
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		LastName:     req.LastName,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Zip:          req.Zip,
		City:         req.City,
		Region:       req.Region,
		Country:      req.Country,
		TaxNumber:    req.TaxNumber,
		BankName:     req.BankName,
		BankAddress:  req.BankAddress,
		BankCountry:  req.BankCountry,
		BankIBAN:     req.BankIBAN,
		BankBIC:      req.BankBIC,
	}

	consultant, err := h.service.CreateConsultant(c.Request.Context(), input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, toResponse(consultant))
}

func (h *Handler) UpdateConsultant(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateConsultantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeUpdateConsultantRequest(&req)

	input := service.UpdateConsultantInput{
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		LastName:     req.LastName,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Zip:          req.Zip,
		City:         req.City,
		Region:       req.Region,
		Country:      req.Country,
		TaxNumber:    req.TaxNumber,
		BankName:     req.BankName,
		BankAddress:  req.BankAddress,
		BankCountry:  req.BankCountry,
		BankIBAN:     req.BankIBAN,
		BankBIC:      req.BankBIC,
	}

	consultant, err := h.service.UpdateConsultant(c.Request.Context(), id, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(consultant))
}

func (h *Handler) DeleteConsultant(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteConsultant(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetConsultant(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	consultant, err := h.service.GetConsultant(c.Request.Context(), id)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(consultant))
}
