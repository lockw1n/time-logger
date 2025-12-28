package service

import "errors"

var (
	ErrContractNotFound      = errors.New("contract not found")
	ErrContractAlreadyExists = errors.New("contract already exists")
	ErrContractConflict      = errors.New("contract conflict")
	ErrContractInvalid       = errors.New("invalid contract data")
)
