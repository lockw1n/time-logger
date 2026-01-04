package migration

import (
	activityrepo "github.com/lockw1n/time-logger/internal/activity/repository"
	companyrepo "github.com/lockw1n/time-logger/internal/company/repository"
	consultantmigration "github.com/lockw1n/time-logger/internal/consultant/migration"
	consultantrepo "github.com/lockw1n/time-logger/internal/consultant/repository"
	contractrepo "github.com/lockw1n/time-logger/internal/contract/repository"
	entryrepo "github.com/lockw1n/time-logger/internal/entry/repository"
	ticketrepo "github.com/lockw1n/time-logger/internal/ticket/repository"
)

func ModelsForMigration() []any {
	var models []any

	models = append(models, companyrepo.Migrations()...)
	models = append(models, consultantrepo.Migrations()...)
	models = append(models, contractrepo.Migrations()...)
	models = append(models, ticketrepo.Migrations()...)
	models = append(models, activityrepo.Migrations()...)
	models = append(models, entryrepo.Migrations()...)

	return models
}

func ExplicitMigrations() []Migration {
	return []Migration{
		consultantmigration.AddAuthFields{},
	}
}
