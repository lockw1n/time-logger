package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/constants"
	"github.com/lockw1n/time-logger/internal/contract/domain"
)

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func normalizeCreateContractRequest(req *CreateContractRequest) {
	req.Currency = strings.TrimSpace(req.Currency)
	req.OrderNumber = strings.TrimSpace(req.OrderNumber)
	req.StartDate = strings.TrimSpace(req.StartDate)

	req.EndDate = trimPtr(req.EndDate)
	req.PaymentTerms = trimPtr(req.PaymentTerms)
}

func normalizeUpdateContractRequest(req *UpdateContractRequest) {
	req.Currency = trimPtr(req.Currency)
	req.OrderNumber = trimPtr(req.OrderNumber)
	req.PaymentTerms = trimPtr(req.PaymentTerms)
	req.StartDate = trimPtr(req.StartDate)
	req.EndDate = trimPtr(req.EndDate)
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

func parseDatePtr(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	t, err := parseDate(*value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(constants.InternalDateFormat, date)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return t, nil
}

func respondContracts(c *gin.Context, contracts []domain.Contract) {
	resp := make([]ContractResponse, 0, len(contracts))
	for _, contract := range contracts {
		resp = append(resp, toResponse(contract))
	}
	c.JSON(http.StatusOK, resp)
}
