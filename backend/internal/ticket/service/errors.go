package service

import "errors"

var (
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrTicketAlreadyExists = errors.New("ticket already exists")
	ErrTicketConflict      = errors.New("ticket conflict")
	ErrTicketInvalid       = errors.New("invalid ticket data")
)
