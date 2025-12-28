package service

import "errors"

var (
	ErrEntryNotFound      = errors.New("entry not found")
	ErrEntryAlreadyExists = errors.New("entry already exists")
	ErrEntryConflict      = errors.New("entry conflict")
	ErrEntryInvalid       = errors.New("invalid entry data")
)
