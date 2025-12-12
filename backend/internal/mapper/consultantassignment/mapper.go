package consultantassignment

import (
	"time"

	assignmentdto "github.com/lockw1n/time-logger/internal/dto/consultantassignment"
	"github.com/lockw1n/time-logger/internal/models"
)

/*************** DTO → Model ***************/

func ToModel(d assignmentdto.Request) (*models.ConsultantAssignment, error) {
	var startDate, endDate *time.Time
	var err error

	if d.StartDate != nil {
		t, err := time.Parse("2006-01-02", *d.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &t
	}

	if d.EndDate != nil {
		t, err := time.Parse("2006-01-02", *d.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &t
	}

	return &models.ConsultantAssignment{
		ConsultantID: d.ConsultantID,
		CompanyID:    d.CompanyID,
		HourlyRate:   d.HourlyRate,
		Currency:     d.Currency,
		OrderNumber:  d.OrderNumber,
		StartDate:    startDate,
		EndDate:      endDate,
	}, err
}

/*************** Model → DTO ***************/

func ToResponse(m *models.ConsultantAssignment) assignmentdto.Response {
	var start, end *string

	if m.StartDate != nil {
		s := m.StartDate.Format("2006-01-02")
		start = &s
	}

	if m.EndDate != nil {
		s := m.EndDate.Format("2006-01-02")
		end = &s
	}

	return assignmentdto.Response{
		ID:           m.ID,
		ConsultantID: m.ConsultantID,
		CompanyID:    m.CompanyID,
		HourlyRate:   m.HourlyRate,
		Currency:     m.Currency,
		OrderNumber:  m.OrderNumber,
		StartDate:    start,
		EndDate:      end,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}
