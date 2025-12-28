package excel

type InvoiceView struct {
	Number     string
	Activities []ActivityView
}

type ActivityView struct {
	Title   string
	Entries []EntryView
}

type EntryView struct {
	DateFormatted string
	TicketCode    string
	Hours         float64
}
