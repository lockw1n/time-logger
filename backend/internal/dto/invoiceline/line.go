package invoiceline

type Line struct {
	ID      uint64  `json:"id"`
	EntryID uint64  `json:"entry_id"`
	Hours   float64 `json:"hours"`
	Rate    float64 `json:"rate"`
	Amount  float64 `json:"amount"`

	CreatedAt string `json:"created_at"`
}
