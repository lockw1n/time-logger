package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/company/service"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrCompanyNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrCompanyInvalid):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrCompanyAlreadyExists),
		errors.Is(err, service.ErrCompanyConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
