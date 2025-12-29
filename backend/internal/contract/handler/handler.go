package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/contract/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateContract(c *gin.Context) {
	var req CreateContractRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeCreateContractRequest(&req)

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	endDate, err := parseDatePtr(req.EndDate)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	input := service.CreateContractInput{
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
		HourlyRate:   req.HourlyRate,
		Currency:     req.Currency,
		OrderNumber:  req.OrderNumber,
		PaymentTerms: req.PaymentTerms,
		StartDate:    startDate,
		EndDate:      endDate,
	}

	contract, err := h.service.CreateContract(c.Request.Context(), input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, toResponse(contract))
}

func (h *Handler) UpdateContract(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeUpdateContractRequest(&req)

	startDate, err := parseDatePtr(req.StartDate)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	endDate, err := parseDatePtr(req.EndDate)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	input := service.UpdateContractInput{
		HourlyRate:   req.HourlyRate,
		Currency:     req.Currency,
		OrderNumber:  req.OrderNumber,
		PaymentTerms: req.PaymentTerms,
		StartDate:    startDate,
		EndDate:      endDate,
	}

	contract, err := h.service.UpdateContract(c.Request.Context(), id, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(contract))
}

func (h *Handler) DeleteContract(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteContract(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetContract(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	contract, err := h.service.GetContract(c.Request.Context(), id)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(contract))
}

func (h *Handler) ListContractsForConsultant(c *gin.Context) {
	consultantIDStr := c.Query("consultant_id")
	if consultantIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "consultant_id is required"})
		return
	}

	consultantID, err := strconv.ParseUint(consultantIDStr, 10, 64)
	if err != nil || consultantID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid consultant_id"})
		return
	}

	contracts, err := h.service.ListContractsForConsultant(c.Request.Context(), consultantID)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	respondContracts(c, contracts)
}
