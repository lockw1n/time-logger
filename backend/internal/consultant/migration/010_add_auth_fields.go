package migration

import "gorm.io/gorm"

type AddAuthFields struct{}

func (AddAuthFields) Name() string {
	return "consultants:010_add_auth_fields"
}

func (AddAuthFields) Run(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		const defaultHash = "$2a$10$N9qo8uLOickgx2ZMRZo4i.eWJk2zF4Y4Q6ZJxP9YxV9r7k4dQwYq6" // "change-me"

		if err := tx.Exec(`
			UPDATE consultants
			SET
				email = CONCAT('consultant_', id, '@local'),
				password_hash = ?
			WHERE email IS NULL
		`, defaultHash).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE consultants
				ALTER COLUMN email SET NOT NULL,
				ALTER COLUMN password_hash SET NOT NULL
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS ux_consultants_email
			ON consultants(email)
		`).Error; err != nil {
			return err
		}

		return nil
	})
}
