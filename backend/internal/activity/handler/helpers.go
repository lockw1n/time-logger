package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/activity/domain"
)

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func normalizeCreateActivityRequest(req *CreateActivityRequest) {
	req.Name = strings.TrimSpace(req.Name)
	req.Color = trimPtr(req.Color)
}

func normalizeUpdateActivityRequest(req *UpdateActivityRequest) {
	req.Name = trimPtr(req.Name)
	req.Color = trimPtr(req.Color)
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

func respondActivities(c *gin.Context, activities []domain.Activity) {
	resp := make([]ActivityResponse, 0, len(activities))
	for _, activity := range activities {
		resp = append(resp, toResponse(activity))
	}
	c.JSON(http.StatusOK, resp)
}
