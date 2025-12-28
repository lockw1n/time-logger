package repository

import "github.com/lockw1n/time-logger/internal/entry/domain"

func toModel(d domain.Entry) entryModel {
	return entryModel{
		ID:              d.ID,
		ConsultantID:    d.ConsultantID,
		CompanyID:       d.CompanyID,
		ContractID:      d.ContractID,
		TicketID:        d.TicketID,
		ActivityID:      d.ActivityID,
		Date:            d.Date,
		DurationMinutes: d.DurationMinutes,
		Comment:         d.Comment,
	}
}

func toDomain(m entryModel) domain.Entry {
	return domain.Entry{
		ID:              m.ID,
		ConsultantID:    m.ConsultantID,
		CompanyID:       m.CompanyID,
		ContractID:      m.ContractID,
		TicketID:        m.TicketID,
		ActivityID:      m.ActivityID,
		Date:            m.Date,
		DurationMinutes: m.DurationMinutes,
		Comment:         m.Comment,
	}
}

func toDomainSlice(models []entryModel) []domain.Entry {
	out := make([]domain.Entry, 0, len(models))
	for _, m := range models {
		out = append(out, toDomain(m))
	}
	return out
}
