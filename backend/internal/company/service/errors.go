package service

import "errors"

var (
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyAlreadyExists = errors.New("company already exists")
	ErrCompanyConflict      = errors.New("company conflict")
	ErrCompanyInvalid       = errors.New("invalid company data")
)
