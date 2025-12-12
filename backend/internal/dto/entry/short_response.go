package entry

type ShortResponse struct {
	ID              uint64  `json:"id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}
