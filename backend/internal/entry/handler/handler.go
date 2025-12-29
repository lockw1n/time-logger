package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/entry/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateEntry(c *gin.Context) {
	var req CreateEntryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeCreateEntryRequest(&req)

	date, err := parseDate(req.Date)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	input := service.CreateEntryInput{
		ConsultantID:    req.ConsultantID,
		CompanyID:       req.CompanyID,
		TicketCode:      req.TicketCode,
		ActivityID:      req.ActivityID,
		Date:            date,
		DurationMinutes: req.DurationMinutes,
		Comment:         req.Comment,
	}

	entry, err := h.service.CreateEntry(c.Request.Context(), input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, toResponse(entry))
}

func (h *Handler) UpdateEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	normalizeUpdateEntryRequest(&req)

	input := service.UpdateEntryInput{
		DurationMinutes: req.DurationMinutes,
		Comment:         req.Comment,
	}

	entry, err := h.service.UpdateEntry(c.Request.Context(), id, input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(entry))
}

func (h *Handler) DeleteEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteEntry(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.service.GetEntry(c.Request.Context(), id)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(entry))
}
