package repository

import "github.com/lockw1n/time-logger/internal/activity/domain"

func toModel(d domain.Activity) activityModel {
	return activityModel{
		ID:        d.ID,
		CompanyID: d.CompanyID,
		Name:      d.Name,
		Color:     d.Color,
		Billable:  d.Billable,
		Priority:  d.Priority,
	}
}

func toDomain(m activityModel) domain.Activity {
	return domain.Activity{
		ID:        m.ID,
		CompanyID: m.CompanyID,
		Name:      m.Name,
		Color:     m.Color,
		Billable:  m.Billable,
		Priority:  m.Priority,
	}
}

func toDomainSlice(models []activityModel) []domain.Activity {
	out := make([]domain.Activity, 0, len(models))
	for _, m := range models {
		out = append(out, toDomain(m))
	}
	return out
}
