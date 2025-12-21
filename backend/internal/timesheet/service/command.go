package service

type GenerateReportCommand struct {
	ConsultantID uint64
	CompanyID    uint64
	Start        string
	End          string
}
