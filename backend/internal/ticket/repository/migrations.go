package repository

func Migrations() []any {
	return []any{
		&ticketModel{},
	}
}
