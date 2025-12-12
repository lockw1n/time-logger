package entry

type Response struct {
	ID                     uint64  `json:"id"`
	ConsultantID           uint64  `json:"consultant_id"`
	CompanyID              uint64  `json:"company_id"`
	ConsultantAssignmentID uint64  `json:"consultant_assignment_id"`
	TicketID               *uint64 `json:"ticket_id"`
	LabelID                *uint64 `json:"label_id"`

	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`

	TicketLabel *string `json:"ticket_label,omitempty"`
	LabelName   *string `json:"label_name,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
