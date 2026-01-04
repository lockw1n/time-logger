package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authctx "github.com/lockw1n/time-logger/internal/auth/context"
	"github.com/lockw1n/time-logger/internal/timesheet/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTimesheet(c *gin.Context) {
	auth, ok := authctx.FromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth context missing"})
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
		ConsultantID: auth.ConsultantID,
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
