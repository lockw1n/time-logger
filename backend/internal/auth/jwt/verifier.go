package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Verifier interface {
	VerifyAccessToken(tokenString string) (Claims, error)
}

type verifier struct {
	config Config
	key    []byte
}

func NewVerifier(config Config) Verifier {
	return &verifier{
		config: config,
		key:    []byte(config.Secret),
	}
}

func (v *verifier) VerifyAccessToken(tokenString string) (Claims, error) {
	var claims Claims

	parser := jwtlib.NewParser(
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
		jwtlib.WithAudience(v.config.Audience),
		jwtlib.WithIssuer(v.config.Issuer),
		jwtlib.WithLeeway(30*time.Second),
	)

	_, err := parser.ParseWithClaims(tokenString, &claims, func(token *jwtlib.Token) (any, error) {
		return v.key, nil
	})

	if err == nil {
		return claims, nil
	}

	if errors.Is(err, jwtlib.ErrTokenExpired) {
		return Claims{}, ErrExpiredToken
	}

	return Claims{}, ErrInvalidToken
}
