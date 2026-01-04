package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authctx "github.com/lockw1n/time-logger/internal/auth/context"
	authjwt "github.com/lockw1n/time-logger/internal/auth/jwt"
)

type Middleware struct {
	verifier authjwt.Verifier
}

func NewMiddleware(verifier authjwt.Verifier) *Middleware {
	return &Middleware{verifier: verifier}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header missing",
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		claims, err := m.verifier.VerifyAccessToken(token)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, authjwt.ErrExpiredToken) {
				status = http.StatusUnauthorized
			}

			c.AbortWithStatusJSON(status, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		ctx := authctx.WithAuth(
			c.Request.Context(),
			authctx.AuthContext{
				ConsultantID: claims.ConsultantID,
			},
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
