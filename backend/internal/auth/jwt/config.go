package jwt

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Secret   string
	Issuer   string
	Audience string
	TTL      time.Duration
}

func LoadConfigFromEnv() (Config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "time-logger"
	}

	audience := os.Getenv("JWT_AUDIENCE")
	if audience == "" {
		audience = "time-logger-clients"
	}

	ttlStr := os.Getenv("JWT_TTL")
	ttl := 30 * time.Minute // sensible default
	if ttlStr != "" {
		parsed, err := time.ParseDuration(ttlStr)
		if err != nil {
			return Config{}, err
		}
		ttl = parsed
	}

	return Config{
		Secret:   secret,
		Issuer:   issuer,
		Audience: audience,
		TTL:      ttl,
	}, nil
}
