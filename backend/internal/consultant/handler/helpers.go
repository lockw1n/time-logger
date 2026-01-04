package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	authctx "github.com/lockw1n/time-logger/internal/auth/context"
)

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func normalizeCreateConsultantRequest(req *CreateConsultantRequest) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.AddressLine1 = strings.TrimSpace(req.AddressLine1)
	req.Zip = strings.TrimSpace(req.Zip)
	req.City = strings.TrimSpace(req.City)
	req.Country = strings.TrimSpace(req.Country)
	req.TaxNumber = strings.TrimSpace(req.TaxNumber)
	req.BankName = strings.TrimSpace(req.BankName)
	req.BankAddress = strings.TrimSpace(req.BankAddress)
	req.BankCountry = strings.TrimSpace(req.BankCountry)
	req.BankIBAN = strings.TrimSpace(req.BankIBAN)
	req.BankBIC = strings.TrimSpace(req.BankBIC)

	req.MiddleName = trimPtr(req.MiddleName)
	req.AddressLine2 = trimPtr(req.AddressLine2)
	req.Region = trimPtr(req.Region)
}

func normalizeUpdateConsultantRequest(req *UpdateConsultantRequest) {
	req.FirstName = trimPtr(req.FirstName)
	req.MiddleName = trimPtr(req.MiddleName)
	req.LastName = trimPtr(req.LastName)
	req.AddressLine1 = trimPtr(req.AddressLine1)
	req.AddressLine2 = trimPtr(req.AddressLine2)
	req.Zip = trimPtr(req.Zip)
	req.City = trimPtr(req.City)
	req.Region = trimPtr(req.Region)
	req.Country = trimPtr(req.Country)
	req.TaxNumber = trimPtr(req.TaxNumber)
	req.BankName = trimPtr(req.BankName)
	req.BankAddress = trimPtr(req.BankAddress)
	req.BankCountry = trimPtr(req.BankCountry)
	req.BankIBAN = trimPtr(req.BankIBAN)
	req.BankBIC = trimPtr(req.BankBIC)
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

func requireConsultantID(c *gin.Context) (uint64, bool) {
	auth, ok := authctx.FromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "auth context missing",
		})
		return 0, false
	}

	return auth.ConsultantID, true
}
