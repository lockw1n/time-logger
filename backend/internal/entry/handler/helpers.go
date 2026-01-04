package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	authctx "github.com/lockw1n/time-logger/internal/auth/context"
	"github.com/lockw1n/time-logger/internal/constants"
)

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func normalizeCreateEntryRequest(req *CreateEntryRequest) {
	req.Date = strings.TrimSpace(req.Date)
	req.Comment = trimPtr(req.Comment)
}

func normalizeUpdateEntryRequest(req *UpdateEntryRequest) {
	req.Comment = trimPtr(req.Comment)
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(constants.InternalDateFormat, date)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return t, nil
}

func requireConsultantID(c *gin.Context) (uint64, bool) {
	auth, ok := authctx.FromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "auth context missing",
		})
		return 0, false
	}

	return auth.ConsultantID, true
}
