package service

import "errors"

var (
	ErrActivityNotFound      = errors.New("activity not found")
	ErrActivityAlreadyExists = errors.New("activity already exists")
	ErrActivityConflict      = errors.New("activity conflict")
	ErrActivityInvalid       = errors.New("invalid activity data")
)
