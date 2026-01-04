package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	authjwt "github.com/lockw1n/time-logger/internal/auth/jwt"
	consultantrepo "github.com/lockw1n/time-logger/internal/consultant/repository"
)

type service struct {
	consultantRepo consultantrepo.Repository
	issuer         authjwt.Issuer
}

func NewService(
	consultantRepo consultantrepo.Repository,
	issuer authjwt.Issuer,
) Service {
	return &service{
		consultantRepo: consultantRepo,
		issuer:         issuer,
	}
}

func (s *service) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	password := strings.TrimSpace(in.Password)

	if email == "" || password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	consultant, err := s.consultantRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, consultantrepo.ErrNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(consultant.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, expiresAt, err := s.issuer.IssueAccessToken(consultant.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		ConsultantID: consultant.ID,
		AccessToken:  token,
		ExpiresAt:    expiresAt,
	}, nil
}
