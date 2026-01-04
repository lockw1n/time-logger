package handler

import (
	"errors"
	"net/http"

	authservice "github.com/lockw1n/time-logger/internal/auth/service"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, authservice.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
