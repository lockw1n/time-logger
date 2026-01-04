package migration

import "gorm.io/gorm"

type Migration interface {
	Name() string
	Run(db *gorm.DB) error
}
