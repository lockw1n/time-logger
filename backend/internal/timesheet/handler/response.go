package handler

import (
	activitydomain "github.com/lockw1n/time-logger/internal/activity/domain"
	"github.com/lockw1n/time-logger/internal/constants"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type TimesheetResponse struct {
	ConsultantID uint64                  `json:"consultant_id"`
	CompanyID    uint64                  `json:"company_id"`
	Start        string                  `json:"start"`
	End          string                  `json:"end"`
	Rows         []TimesheetRowResponse  `json:"rows"`
	Totals       TimesheetTotalsResponse `json:"totals"`
}

type TimesheetRowResponse struct {
	Ticket        TimesheetTicketResponse   `json:"ticket"`
	Activity      TimesheetActivityResponse `json:"activity"`
	Entries       []TimesheetEntryResponse  `json:"entries"`
	PerDayMinutes map[string]int            `json:"per_day_minutes"`
	TotalMinutes  int                       `json:"total_minutes"`
}

type TimesheetTicketResponse struct {
	ID          uint64  `json:"id"`
	CompanyID   uint64  `json:"company_id"`
	Code        string  `json:"code"`
	Title       *string `json:"title"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

type TimesheetActivityResponse struct {
	ID        uint64  `json:"id"`
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
	Billable  bool    `json:"billable"`
	Priority  int     `json:"priority"`
}

type TimesheetEntryResponse struct {
	ID              uint64  `json:"id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}

type TimesheetTotalsResponse struct {
	PerDayMinutes  map[string]int `json:"per_day_minutes"`
	OverallMinutes int            `json:"overall_minutes"`
}

func toResponse(timesheet domain.Timesheet) TimesheetResponse {
	return TimesheetResponse{
		ConsultantID: timesheet.ConsultantID,
		CompanyID:    timesheet.CompanyID,
		Start:        timesheet.Start.Format(constants.ResponseDateFormat),
		End:          timesheet.End.Format(constants.ResponseDateFormat),
		Rows:         toRowResponse(timesheet.Rows),
		Totals:       toTotalsResponse(timesheet.Totals),
	}
}

func toRowResponse(rows []domain.TimesheetRow) []TimesheetRowResponse {
	resp := make([]TimesheetRowResponse, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, TimesheetRowResponse{
			Ticket:        toTicketResponse(row.Ticket),
			Activity:      toActivityResponse(row.Activity),
			Entries:       toEntryResponse(row.Entries),
			PerDayMinutes: row.PerDayMinutes,
			TotalMinutes:  row.TotalMinutes,
		})
	}
	return resp
}

func toTicketResponse(ticket ticketdomain.Ticket) TimesheetTicketResponse {
	return TimesheetTicketResponse{
		ID:          ticket.ID,
		CompanyID:   ticket.CompanyID,
		Code:        ticket.Code,
		Title:       ticket.Title,
		Label:       ticket.Label,
		Description: ticket.Description,
	}
}

func toActivityResponse(activity activitydomain.Activity) TimesheetActivityResponse {
	return TimesheetActivityResponse{
		ID:        activity.ID,
		CompanyID: activity.CompanyID,
		Name:      activity.Name,
		Color:     activity.Color,
		Billable:  activity.Billable,
		Priority:  activity.Priority,
	}
}

func toEntryResponse(entries []domain.TimesheetEntry) []TimesheetEntryResponse {
	resp := make([]TimesheetEntryResponse, 0, len(entries))
	for _, entry := range entries {
		resp = append(resp, TimesheetEntryResponse{
			ID:              entry.ID,
			Date:            entry.Date.Format(constants.ResponseDateFormat),
			DurationMinutes: entry.DurationMinutes,
			Comment:         entry.Comment,
		})
	}
	return resp
}

func toTotalsResponse(totals domain.TimesheetTotals) TimesheetTotalsResponse {
	return TimesheetTotalsResponse{
		PerDayMinutes:  totals.PerDayMinutes,
		OverallMinutes: totals.OverallMinutes,
	}
}
