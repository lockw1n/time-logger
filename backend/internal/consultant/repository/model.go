package repository

import "time"

type consultantModel struct {
	ID         uint64  `gorm:"primaryKey;autoIncrement"`
	FirstName  string  `gorm:"column:first_name;not null"`
	MiddleName *string `gorm:"column:middle_name"`
	LastName   string  `gorm:"column:last_name;not null"`

	Email        string `gorm:"column:email;size:255"`
	PasswordHash string `gorm:"column:password_hash;size:255"`

	AddressLine1 string  `gorm:"column:address_line1;not null"`
	AddressLine2 *string `gorm:"column:address_line2"`
	Zip          string  `gorm:"column:zip;not null"`
	City         string  `gorm:"column:city;not null"`
	Region       *string `gorm:"column:region"`
	Country      string  `gorm:"column:country;not null"`

	TaxNumber string `gorm:"column:tax_number; not null"`

	BankName    string `gorm:"column:bank_name;not null"`
	BankAddress string `gorm:"column:bank_address;not null"`
	BankCountry string `gorm:"column:bank_country;not null"`
	BankIBAN    string `gorm:"column:bank_iban;not null"`
	BankBIC     string `gorm:"column:bank_bic;not null"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (consultantModel) TableName() string {
	return "consultants"
}
