package handler

import (
	"github.com/lockw1n/time-logger/internal/constants"
	"github.com/lockw1n/time-logger/internal/entry/domain"
)

type EntryResponse struct {
	ID              uint64  `json:"id"`
	ConsultantID    uint64  `json:"consultant_id"`
	CompanyID       uint64  `json:"company_id"`
	ContractID      uint64  `json:"contract_id"`
	TicketID        uint64  `json:"ticket_id"`
	ActivityID      uint64  `json:"activity_id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}

func toResponse(entry domain.Entry) EntryResponse {
	return EntryResponse{
		ID:              entry.ID,
		ConsultantID:    entry.ConsultantID,
		CompanyID:       entry.CompanyID,
		ContractID:      entry.ContractID,
		TicketID:        entry.TicketID,
		ActivityID:      entry.ActivityID,
		Date:            entry.Date.Format(constants.ResponseDateFormat),
		DurationMinutes: entry.DurationMinutes,
		Comment:         entry.Comment,
	}
}
