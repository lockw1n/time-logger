package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/contract/service"
)

var (
	ErrInvalidDateFormat = errors.New("date must be in YYYY-MM-DD format")
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrContractNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrContractInvalid),
		errors.Is(err, ErrInvalidDateFormat):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrContractAlreadyExists),
		errors.Is(err, service.ErrContractConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
