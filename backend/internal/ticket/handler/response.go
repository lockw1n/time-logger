package handler

import (
	"github.com/lockw1n/time-logger/internal/ticket/domain"
)

type TicketResponse struct {
	ID          uint64  `json:"id"`
	CompanyID   uint64  `json:"company_id"`
	Code        string  `json:"code"`
	Title       *string `json:"title"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

func toResponse(ticket domain.Ticket) TicketResponse {
	return TicketResponse{
		ID:          ticket.ID,
		CompanyID:   ticket.CompanyID,
		Code:        ticket.Code,
		Title:       ticket.Title,
		Label:       ticket.Label,
		Description: ticket.Description,
	}
}
