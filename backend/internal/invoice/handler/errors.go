package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/invoice/service"
)

var (
	ErrInvalidDateFormat = errors.New("date must be in YYYY-MM-DD format")
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvoiceInvalid),
		errors.Is(err, ErrInvalidDateFormat):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrInvoiceConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
