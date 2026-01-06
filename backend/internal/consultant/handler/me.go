package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/consultant/service"
)

func (h *Handler) GetMe(c *gin.Context) {
	consultantID, ok := requireConsultantID(c)
	if !ok {
		return
	}

	consultant, err := h.service.GetMe(c.Request.Context(), consultantID)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(consultant))
}

func (h *Handler) UpdateMe(c *gin.Context) {
	consultantID, ok := requireConsultantID(c)
	if !ok {
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
		Email:        req.Email,
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

	consultant, err := h.service.UpdateMe(
		c.Request.Context(),
		consultantID,
		input,
	)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(consultant))
}
