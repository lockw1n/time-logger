package service

import (
	"strings"
)

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i CreateConsultantInput) Validate() error {
	if strings.TrimSpace(i.FirstName) == "" {
		return validationError("firstname should not be empty")
	}
	if strings.TrimSpace(i.LastName) == "" {
		return validationError("lastname should not be empty")
	}

	return nil
}

func (i UpdateConsultantInput) Validate() error {
	if i.FirstName != nil && strings.TrimSpace(*i.FirstName) == "" {
		return validationError("firstname should not be empty")
	}
	if i.LastName != nil && strings.TrimSpace(*i.LastName) == "" {
		return validationError("lastname should not be empty")
	}

	return nil
}
