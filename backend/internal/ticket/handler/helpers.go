package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/ticket/domain"
)

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func normalizeCreateTicketRequest(req *CreateTicketRequest) {
	req.Code = strings.TrimSpace(req.Code)
	req.Title = trimPtr(req.Title)
	req.Label = trimPtr(req.Label)
	req.Description = trimPtr(req.Description)
}

func normalizeUpdateTicketRequest(req *UpdateTicketRequest) {
	req.Title = trimPtr(req.Title)
	req.Label = trimPtr(req.Label)
	req.Description = trimPtr(req.Description)
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func respondTickets(c *gin.Context, tickets []domain.Ticket) {
	resp := make([]TicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		resp = append(resp, toResponse(ticket))
	}
	c.JSON(http.StatusOK, resp)
}
