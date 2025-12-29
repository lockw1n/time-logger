package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/company/domain"
)

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func normalizeCreateCompanyRequest(req *CreateCompanyRequest) {
	req.Name = strings.TrimSpace(req.Name)
	req.TaxNumber = strings.TrimSpace(req.TaxNumber)
	req.AddressLine1 = strings.TrimSpace(req.AddressLine1)
	req.Zip = strings.TrimSpace(req.Zip)
	req.City = strings.TrimSpace(req.City)
	req.Country = strings.TrimSpace(req.Country)

	req.NameShort = trimPtr(req.NameShort)
	req.AddressLine2 = trimPtr(req.AddressLine2)
	req.Region = trimPtr(req.Region)
}

func normalizeUpdateCompanyRequest(req *UpdateCompanyRequest) {
	req.Name = trimPtr(req.Name)
	req.NameShort = trimPtr(req.NameShort)
	req.TaxNumber = trimPtr(req.TaxNumber)
	req.AddressLine1 = trimPtr(req.AddressLine1)
	req.AddressLine2 = trimPtr(req.AddressLine2)
	req.Zip = trimPtr(req.Zip)
	req.City = trimPtr(req.City)
	req.Region = trimPtr(req.Region)
	req.Country = trimPtr(req.Country)
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

func respondCompanies(c *gin.Context, companies []domain.Company) {
	resp := make([]CompanyResponse, 0, len(companies))
	for _, company := range companies {
		resp = append(resp, toResponse(company))
	}
	c.JSON(http.StatusOK, resp)
}
