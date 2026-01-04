package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Issuer interface {
	IssueAccessToken(consultantID uint64) (token string, expiresAt time.Time, err error)
}

type issuer struct {
	cfg Config
	key []byte
}

func NewIssuer(cfg Config) Issuer {
	return &issuer{
		cfg: cfg,
		key: []byte(cfg.Secret),
	}
}

type Claims struct {
	ConsultantID uint64 `json:"consultant_id"`
	jwtlib.RegisteredClaims
}

func (i *issuer) IssueAccessToken(consultantID uint64) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(i.cfg.TTL)

	claims := Claims{
		ConsultantID: consultantID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    i.cfg.Issuer,
			Subject:   "consultant",
			Audience:  jwtlib.ClaimStrings{i.cfg.Audience},
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	signed, err := token.SignedString(i.key)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}
