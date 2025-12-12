package timesheet

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lockw1n/time-logger/internal/constants"
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

	// Bind query parameters into DTO
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	// Validate required numeric IDs
	if req.ConsultantID == 0 || req.CompanyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "consultant_id and company_id are required"})
		return
	}

	// Validate dates exist
	if req.Start == "" || req.End == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end dates are required"})
		return
	}

	// Validate date format
	if _, err := time.Parse(constants.DateFormat, req.Start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}

	if _, err := time.Parse(constants.DateFormat, req.End); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}

	// Generate report
	report, err := h.service.GenerateReport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate timesheet"})
		return
	}

	c.JSON(http.StatusOK, report)
}
