package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authservice "github.com/lockw1n/time-logger/internal/auth/service"
)

type Handler struct {
	service authservice.Service
}

func NewHandler(service authservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.Login(c.Request.Context(), authservice.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: result.AccessToken,
		ExpiresAt:   result.ExpiresAt,
	})
}
