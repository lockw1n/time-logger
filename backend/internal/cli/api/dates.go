package api

import "time"

// The backend is asymmetric about dates: it accepts YYYY-MM-DD in requests
// (constants.InternalDateFormat) but formats scalar date fields — start, end,
// and entry dates — as DD.MM.YYYY in responses (constants.ResponseDateFormat).
// Map keys such as per_day_minutes are already YYYY-MM-DD. To spare the rest of
// the CLI from that quirk, the client normalizes every response date field to
// canonical YYYY-MM-DD right after decoding, so the whole tool works in one
// format (and matches what the CLI sends).
const (
	responseDateFormat  = "02.01.2006"
	canonicalDateFormat = "2006-01-02"
)

// canonicalDate converts a DD.MM.YYYY response date to YYYY-MM-DD. Values that
// are empty or already canonical (or otherwise unparseable in the response
// format) are returned unchanged.
func canonicalDate(s string) string {
	t, err := time.Parse(responseDateFormat, s)
	if err != nil {
		return s
	}
	return t.Format(canonicalDateFormat)
}

// normalizeTimesheetDates rewrites the scalar date fields of a timesheet
// response in place to canonical YYYY-MM-DD.
func normalizeTimesheetDates(ts *TimesheetResponse) {
	ts.Start = canonicalDate(ts.Start)
	ts.End = canonicalDate(ts.End)
	for i := range ts.Rows {
		for j := range ts.Rows[i].Entries {
			ts.Rows[i].Entries[j].Date = canonicalDate(ts.Rows[i].Entries[j].Date)
		}
	}
}
