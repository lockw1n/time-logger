package domain

type Group struct {
	Label string
	Rows  []Row

	TotalHours float64
	HourlyRate float64
	Subtotal   int64 // smallest currency unit (e.g. cents)
}
