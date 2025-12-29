package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/entry/service"
)

var (
	ErrInvalidDateFormat = errors.New("date must be in YYYY-MM-DD format")
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrEntryNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrEntryInvalid),
		errors.Is(err, ErrInvalidDateFormat):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrEntryAlreadyExists),
		errors.Is(err, service.ErrEntryConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
