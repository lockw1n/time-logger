package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/activity/service"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrActivityNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrActivityInvalid):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrActivityAlreadyExists),
		errors.Is(err, service.ErrActivityConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
