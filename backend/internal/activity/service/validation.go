package service

import (
	"strings"
)

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i CreateActivityInput) Validate() error {
	if i.CompanyID == 0 {
		return validationError("company_id is required")
	}
	if i.Name == "" {
		return validationError("name is required")
	}
	if i.Priority < 1 {
		return validationError("priority should be greater than 0")
	}

	return nil
}

func (i UpdateActivityInput) Validate() error {
	if i.Priority != nil && *i.Priority < 1 {
		return validationError("priority should be greater than 0")
	}
	if i.Name != nil && strings.TrimSpace(*i.Name) == "" {
		return validationError("company name should not be empty")
	}

	return nil
}
