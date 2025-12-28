package handler

import (
	"errors"
	"net/http"

	"github.com/lockw1n/time-logger/internal/ticket/service"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrTicketNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrTicketInvalid):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrTicketAlreadyExists),
		errors.Is(err, service.ErrTicketConflict):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
