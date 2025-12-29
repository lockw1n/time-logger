package service

import "time"

type GenerateTimesheetInput struct {
	ConsultantID uint64
	CompanyID    uint64
	Start        time.Time
	End          time.Time
}
