package repository

import "github.com/lockw1n/time-logger/internal/ticket/domain"

func toModel(d domain.Ticket) ticketModel {
	return ticketModel{
		ID:          d.ID,
		CompanyID:   d.CompanyID,
		Code:        d.Code,
		Title:       d.Title,
		Label:       d.Label,
		Description: d.Description,
	}
}

func toDomain(m ticketModel) domain.Ticket {
	return domain.Ticket{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		Code:        m.Code,
		Title:       m.Title,
		Label:       m.Label,
		Description: m.Description,
	}
}

func toDomainSlice(models []ticketModel) []domain.Ticket {
	out := make([]domain.Ticket, 0, len(models))
	for _, m := range models {
		out = append(out, toDomain(m))
	}
	return out
}
