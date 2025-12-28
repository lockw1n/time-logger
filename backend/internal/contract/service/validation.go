package service

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i CreateContractInput) Validate() error {
	if i.ConsultantID == 0 {
		return validationError("consultant_id is required")
	}
	if i.CompanyID == 0 {
		return validationError("company_id is required")
	}
	if i.HourlyRate <= 0 {
		return validationError("hourly_rate is required and must be greater than 0")
	}
	if i.Currency == "" {
		return validationError("currency is required")
	}
	if i.OrderNumber == "" {
		return validationError("order_number is required")
	}
	if i.StartDate.IsZero() {
		return validationError("start_date is required")
	}
	if i.EndDate != nil && i.EndDate.Before(i.StartDate) {
		return validationError("end date before start date")
	}

	return nil
}

func (i UpdateContractInput) Validate() error {
	if i.HourlyRate != nil && *i.HourlyRate <= 0 {
		return validationError("hourly_rate must be greater than 0")
	}

	return nil
}
