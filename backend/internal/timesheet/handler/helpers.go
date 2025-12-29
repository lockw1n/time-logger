package handler

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/constants"
)

func parseRequiredUintQuery(c *gin.Context, name string) (uint64, error) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return 0, errors.New(name + " is required")
	}

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid " + name)
	}

	return id, nil
}

func parseRequiredDateQuery(c *gin.Context, name string) (time.Time, error) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return time.Time{}, errors.New(name + " is required")
	}

	return parseDate(value)
}

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(constants.InternalDateFormat, date)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return t, nil
}
