package service

import (
	"strings"
)

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i CreateCompanyInput) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return validationError("company name should not be empty")
	}
	if strings.TrimSpace(i.Country) == "" {
		return validationError("company country should not be empty")
	}

	return nil
}

func (i UpdateCompanyInput) Validate() error {
	if i.Name != nil && strings.TrimSpace(*i.Name) == "" {
		return validationError("company name should not be empty")
	}

	return nil
}
