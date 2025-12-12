package timesheet

type Request struct {
	ConsultantID uint64 `form:"consultant_id" json:"consultant_id"`
	CompanyID    uint64 `form:"company_id" json:"company_id"`
	Start        string `form:"start" json:"start"`
	End          string `form:"end" json:"end"`
}
