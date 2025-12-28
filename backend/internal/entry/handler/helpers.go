package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/constants"
	"github.com/lockw1n/time-logger/internal/entry/domain"
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

func respondContracts(c *gin.Context, entries []domain.Entry) {
	resp := make([]EntryResponse, 0, len(entries))
	for _, entry := range entries {
		resp = append(resp, toResponse(entry))
	}
	c.JSON(http.StatusOK, resp)
}
