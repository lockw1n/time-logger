package service

type validationError string

func (e validationError) Error() string {
	return string(e)
}

func (i GenerateTimesheetInput) Validate() error {
	if i.ConsultantID == 0 {
		return validationError("consultant_id is required")
	}
	if i.CompanyID == 0 {
		return validationError("company_id is required")
	}
	if i.End.Before(i.Start) {
		return validationError("end date before start date")
	}

	return nil
}
