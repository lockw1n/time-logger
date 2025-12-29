package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/ticket/domain"
)

type Service interface {
	CreateTicket(ctx context.Context, input CreateTicketInput) (domain.Ticket, error)
	UpdateTicket(ctx context.Context, id uint64, input UpdateTicketInput) (domain.Ticket, error)
	DeleteTicket(ctx context.Context, id uint64) error

	GetTicket(ctx context.Context, id uint64) (domain.Ticket, error)
	ListTicketsForCompany(ctx context.Context, companyID uint64) ([]domain.Ticket, error)
}
