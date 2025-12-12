package ticket

import (
	ticketdto "github.com/lockw1n/time-logger/internal/dto/ticket"
	"github.com/lockw1n/time-logger/internal/models"
)

/*************** DTO → Model ***************/

func ToModel(d ticketdto.Request) *models.Ticket {
	return &models.Ticket{
		CompanyID:   d.CompanyID,
		Code:        d.Code,
		Label:       d.Label,
		Description: d.Description,
	}
}

/*************** Model → DTO ***************/

func ToResponse(m *models.Ticket) ticketdto.Response {
	if m == nil {
		return ticketdto.Response{
			ID:    0,
			Code:  "(missing ticket)",
			Label: "",
		}
	}

	return ticketdto.Response{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		Code:        m.Code,
		Label:       m.Label,
		Description: m.Description,
	}
}

func ToResponses(list []models.Ticket) []ticketdto.Response {
	out := make([]ticketdto.Response, len(list))
	for i := range list {
		out[i] = ToResponse(&list[i])
	}
	return out
}
