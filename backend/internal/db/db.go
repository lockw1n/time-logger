package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const hoursConstraintName = "entries_hours_quarter_check"
const labelConstraintName = "entries_label_allowed_check"
const dateConstraintName = "entries_date_midnight_check"

// Connect initializes the database connection.
func Connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = database
	return DB
}

// ApplyConstraints adds check constraints for hours, date, and label.
// Safe to run at creation time (fresh DB) and idempotent for restarts.
func ApplyConstraints(database *gorm.DB, labels []string) {
	// Hours constraint
	hoursCheck := fmt.Sprintf(`
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = '%s') THEN
		ALTER TABLE entries
		ADD CONSTRAINT %s CHECK (hours >= 0 AND hours <= 24 AND mod(hours * 4, 1) = 0);
	END IF;
END;
$$;`, hoursConstraintName, hoursConstraintName)
	if err := database.Exec(hoursCheck).Error; err != nil {
		log.Printf("warning: unable to apply quarter-hour check constraint: %v", err)
	}

	// Date constraint (redundant for DATE type, but guards future schema drift)
	dateCheck := fmt.Sprintf(`
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = '%s') THEN
		ALTER TABLE entries
		ADD CONSTRAINT %s CHECK (date = date::date);
	END IF;
END;
$$;`, dateConstraintName, dateConstraintName)
	if err := database.Exec(dateCheck).Error; err != nil {
		log.Printf("warning: unable to apply date-only check constraint: %v", err)
	}

	// Label constraint (rebuilt each start to reflect env)
	if len(labels) == 0 {
		log.Printf("warning: no labels configured; skipping label constraint")
		return
	}

	quoted := make([]string, 0, len(labels))
	for _, l := range labels {
		safe := strings.ReplaceAll(l, "'", "''")
		quoted = append(quoted, fmt.Sprintf("'%s'", safe))
	}
	labelArray := strings.Join(quoted, ",")

	// Normalize label column to text (safe if already text)
	if err := database.Exec(fmt.Sprintf(`ALTER TABLE entries DROP CONSTRAINT IF EXISTS %s;`, labelConstraintName)).Error; err != nil {
		log.Printf("warning: unable to drop existing label constraint: %v", err)
	}

	constraintSQL := fmt.Sprintf(
		`ALTER TABLE entries ADD CONSTRAINT %s CHECK (label = ANY(ARRAY[%s]::text[]));`,
		labelConstraintName,
		labelArray,
	)
	if err := database.Exec(constraintSQL).Error; err != nil {
		log.Printf("warning: unable to apply label allow-list constraint: %v", err)
	}
}
