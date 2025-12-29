package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/timesheet/service"
)

var (
	ErrInvalidDateFormat = errors.New("date must be in YYYY-MM-DD format")
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrTimesheetInvalid),
		errors.Is(err, ErrInvalidDateFormat):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrTimesheetConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
