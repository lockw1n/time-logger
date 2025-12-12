package timesheet

type Report struct {
	ConsultantID uint64      `json:"consultant_id"`
	CompanyID    uint64      `json:"company_id"`
	Start        string      `json:"start"`
	End          string      `json:"end"`
	Rows         []ReportRow `json:"rows"`
}
