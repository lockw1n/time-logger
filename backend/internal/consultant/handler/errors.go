package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/consultant/service"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrConsultantNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrConsultantInvalid):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrConsultantAlreadyExists),
		errors.Is(err, service.ErrConsultantConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
