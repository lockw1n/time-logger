package service

import "errors"

var (
	ErrInvoiceInvalid  = errors.New("invalid invoice data")
	ErrInvoiceConflict = errors.New("invoice conflict")
)
