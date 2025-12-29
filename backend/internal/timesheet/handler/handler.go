package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/timesheet/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTimesheet(c *gin.Context) {
	consultantID, err := parseRequiredUintQuery(c, "consultant_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, err := parseRequiredUintQuery(c, "company_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start, err := parseRequiredDateQuery(c, "start")
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	end, err := parseRequiredDateQuery(c, "end")
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	input := service.GenerateTimesheetInput{
		ConsultantID: consultantID,
		CompanyID:    companyID,
		Start:        start,
		End:          end,
	}

	timesheet, err := h.service.GenerateTimesheet(c.Request.Context(), input)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, toResponse(timesheet))
}
