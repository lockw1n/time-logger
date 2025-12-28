package handler

import (
	"github.com/lockw1n/time-logger/internal/activity/domain"
)

type ActivityResponse struct {
	ID        uint64  `json:"id"`
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
	Billable  bool    `json:"billable"`
	Priority  int     `json:"priority"`
}

func toResponse(activity domain.Activity) ActivityResponse {
	return ActivityResponse{
		ID:        activity.ID,
		CompanyID: activity.CompanyID,
		Name:      activity.Name,
		Color:     activity.Color,
		Billable:  activity.Billable,
		Priority:  activity.Priority,
	}
}
