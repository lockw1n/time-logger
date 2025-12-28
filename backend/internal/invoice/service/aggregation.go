package service

import (
	"sort"

	activitydomain "github.com/lockw1n/time-logger/internal/activity/domain"
	entrydomain "github.com/lockw1n/time-logger/internal/entry/domain"
	"github.com/lockw1n/time-logger/internal/invoice/domain"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
)

type activityAggregate struct {
	priority int
	activity domain.InvoiceActivity
}

func groupActivities(
	entries []entrydomain.Entry,
	activitiesByID map[uint64]activitydomain.Activity,
	ticketsByID map[uint64]ticketdomain.Ticket,
	hourlyRate float64,
) []domain.InvoiceActivity {

	byActivity := make(map[uint64]*activityAggregate)

	for _, entry := range entries {
		activity, ok := activitiesByID[entry.ActivityID]
		if !ok {
			panic("activity not found for entry")
		}

		ticket, ok := ticketsByID[entry.TicketID]
		if !ok {
			panic("ticket not found for entry")
		}

		agg, exists := byActivity[entry.ActivityID]
		if !exists {
			agg = &activityAggregate{
				priority: activity.Priority,
				activity: domain.InvoiceActivity{
					Name:       activity.Name,
					HourlyRate: hourlyRate,
				},
			}
			byActivity[entry.ActivityID] = agg
		}

		hours := float64(entry.DurationMinutes) / 60.0

		agg.activity.Entries = append(agg.activity.Entries, domain.InvoiceEntry{
			Date:       entry.Date,
			TicketCode: ticket.Code,
			Hours:      hours,
		})

		agg.activity.TotalHours += hours
	}

	// flatten + calculate subtotals
	aggregates := make([]activityAggregate, 0, len(byActivity))
	for _, agg := range byActivity {

		// 🔹 deterministic entry order
		sort.SliceStable(agg.activity.Entries, func(i, j int) bool {
			return agg.activity.Entries[i].Date.Before(agg.activity.Entries[j].Date)
		})

		agg.activity.Subtotal = int64(
			agg.activity.TotalHours * hourlyRate * 100,
		)

		aggregates = append(aggregates, *agg)
	}

	// sort by activity priority (business-defined order)
	sort.SliceStable(aggregates, func(i, j int) bool {
		return aggregates[i].priority < aggregates[j].priority
	})

	// strip helper data
	result := make([]domain.InvoiceActivity, 0, len(aggregates))
	for _, agg := range aggregates {
		result = append(result, agg.activity)
	}

	return result
}

func calculateTotals(activities []domain.InvoiceActivity) domain.InvoiceTotals {
	var totals domain.InvoiceTotals

	for _, a := range activities {
		totals.TotalHours += a.TotalHours
		totals.Subtotal += a.Subtotal
	}

	return totals
}
