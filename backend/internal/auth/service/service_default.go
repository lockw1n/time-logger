package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	consultantrepo "github.com/lockw1n/time-logger/internal/consultant/repository"
)

type service struct {
	consultants consultantrepo.Repository
}

func NewService(consultants consultantrepo.Repository) Service {
	return &service{consultants: consultants}
}

func (s *service) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	password := strings.TrimSpace(in.Password)

	if email == "" || password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	consultant, err := s.consultants.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, consultantrepo.ErrNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(consultant.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	return LoginResult{
		ConsultantID: consultant.ID,
	}, nil
}
