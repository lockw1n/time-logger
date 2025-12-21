package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/timesheet/service"
)

type Timesheet struct {
	service service.Timesheet
}

func NewTimesheet(service service.Timesheet) *Timesheet {
	return &Timesheet{service: service}
}

// GetTimesheet GET /timesheet?consultant_id=1&company_id=2&start=2025-01-01&end=2025-01-14
func (h *Timesheet) GetTimesheet(c *gin.Context) {
	var req TimesheetRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	cmd := service.GenerateReportCommand{
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
		Start:        req.Start,
		End:          req.End,
	}

	report, err := h.service.GenerateReport(c.Request.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingConsultantID),
			errors.Is(err, service.ErrMissingCompanyID),
			errors.Is(err, service.ErrInvalidDateFormat),
			errors.Is(err, service.ErrInvalidDateRange):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate timesheet"})
		}
		return
	}

	c.JSON(http.StatusOK, report)
}
