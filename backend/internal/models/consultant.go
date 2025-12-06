package models

import "time"

// Consultant holds personal/billing details for the consultant.
type Consultant struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	Country     string    `json:"country"`
	Zip         string    `json:"zip"`
	Region      string    `json:"region"`
	City        string    `json:"city"`
	Address1    string    `json:"address_line1"`
	Address2    string    `json:"address_line2"`
	TaxNumber   string    `json:"tax_number"`
	BankName    string    `json:"bank_name"`
	BankAddress string    `json:"bank_address"`
	BankCountry string    `json:"bank_country"`
	IBAN        string    `json:"iban"`
	BIC         string    `json:"bic"`
	OrderNumber string    `json:"order_number"`
	HourlyRate  float64   `json:"hourly_rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
