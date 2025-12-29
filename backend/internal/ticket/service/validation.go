package service

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i CreateTicketInput) Validate() error {
	if i.CompanyID == 0 {
		return validationError("company_id is required")
	}
	if i.Code == "" {
		return validationError("code is required")
	}

	return nil
}
