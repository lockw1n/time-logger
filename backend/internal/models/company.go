package models

import "time"

type Company struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement"`
	Name      string  `gorm:"column:name;not null"`
	NameShort *string `gorm:"column:name_short"`

	TaxNumber *string `gorm:"column:tax_number"`

	AddressLine1 string  `gorm:"column:address_line1;not null"`
	AddressLine2 *string `gorm:"column:address_line2"`
	Zip          string  `gorm:"column:zip;not null"`
	City         string  `gorm:"column:city;not null"`
	Region       *string `gorm:"column:region"`
	Country      string  `gorm:"column:country;not null"`

	PaymentTerms *string `gorm:"column:payment_terms"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Company) TableName() string {
	return "companies"
}
