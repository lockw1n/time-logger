package entry

import (
	"time"

	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
	"github.com/lockw1n/time-logger/internal/models"
)

/*************** DTO → Model ***************/

func ToModel(d entrydto.Request) (*models.Entry, error) {
	date, err := time.Parse("2006-01-02", d.Date)
	if err != nil {
		return nil, err
	}

	return &models.Entry{
		ConsultantID:    d.ConsultantID,
		CompanyID:       d.CompanyID,
		LabelID:         d.LabelID,
		Date:            date,
		DurationMinutes: d.DurationMinutes,
		Comment:         d.Comment,
	}, nil
}

/*************** Model → DTO ***************/

func ToResponse(m *models.Entry) entrydto.Response {
	date := m.Date.Format("2006-01-02")

	resp := entrydto.Response{
		ID:                     m.ID,
		ConsultantID:           m.ConsultantID,
		CompanyID:              m.CompanyID,
		ConsultantAssignmentID: m.ConsultantAssignmentID,
		TicketID:               m.TicketID,
		LabelID:                m.LabelID,
		Date:                   date,
		DurationMinutes:        m.DurationMinutes,
		Comment:                m.Comment,
		CreatedAt:              m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              m.UpdatedAt.Format(time.RFC3339),
	}

	if m.Ticket != nil {
		resp.TicketLabel = &m.Ticket.Label
	}
	if m.Label != nil {
		resp.LabelName = &m.Label.Name
	}

	return resp
}

func ToShortResponse(m *models.Entry) entrydto.ShortResponse {
	return entrydto.ShortResponse{
		ID:              m.ID,
		Date:            m.Date.Format("2006-01-02"),
		DurationMinutes: m.DurationMinutes,
		Comment:         m.Comment,
	}
}
