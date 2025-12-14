package timesheet

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	timesheetdto "github.com/lockw1n/time-logger/internal/dto/timesheet"
	"github.com/lockw1n/time-logger/internal/service/timesheet"
)

type Handler struct {
	service timesheet.Service
}

func NewHandler(service timesheet.Service) *Handler {
	return &Handler{service: service}
}

// GetTimesheet GET /timesheet?consultant_id=1&company_id=2&start=2025-01-01&end=2025-01-14
func (h *Handler) GetTimesheet(c *gin.Context) {
	var req timesheetdto.Request

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	report, err := h.service.GenerateReport(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, timesheet.ErrMissingConsultantID),
			errors.Is(err, timesheet.ErrMissingCompanyID),
			errors.Is(err, timesheet.ErrInvalidDateFormat),
			errors.Is(err, timesheet.ErrInvalidDateRange):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate timesheet"})
		}
		return
	}

	c.JSON(http.StatusOK, report)
}
