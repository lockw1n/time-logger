package service

import "errors"

var (
	ErrConsultantNotFound      = errors.New("consultant not found")
	ErrConsultantAlreadyExists = errors.New("consultant already exists")
	ErrConsultantConflict      = errors.New("consultant conflict")
	ErrConsultantInvalid       = errors.New("invalid consultant data")
)
