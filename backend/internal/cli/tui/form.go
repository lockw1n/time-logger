package tui

import (
	"strings"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// Add-form field indices, in tab order.
const (
	fieldTicket = iota
	fieldActivity
	fieldDuration
	fieldComment
	addFieldCount
)

// addForm is the `a` add-entry form: a ticket code, an activity chosen by cycling
// the fetched list, a duration (phase-2 syntax) and an optional comment. The date
// is fixed to the cursor's day when the form opens. It is a value type so opening
// the form is a pure state transition.
type addForm struct {
	ticket      textInput
	duration    textInput
	comment     textInput
	activities  []api.Activity
	activityIdx int
	field       int
	date        string
}

// newAddForm builds a blank form for the given day. When prefActivity is non-nil
// (enter on an existing row's empty cell) the form is seeded with that row's
// ticket and activity and starts on the duration field. Crucially, the row's
// activity is guaranteed to be the selected one even if it isn't in the fetched
// activities list — it is prepended so the create can never be misfiled under
// activities[0]. A nil prefActivity is the blank `a` form.
func newAddForm(date string, activities []api.Activity, prefTicket string, prefActivity *api.Activity) addForm {
	f := addForm{
		ticket:     newTextInput(prefTicket),
		duration:   newTextInput(""),
		comment:    newTextInput(""),
		activities: activities,
		date:       date,
		field:      fieldTicket,
	}
	if prefActivity != nil {
		// Prefilled from a row: the only thing left to type is the duration.
		f.field = fieldDuration
		if idx := activityIndex(activities, prefActivity.ID); idx >= 0 {
			f.activityIdx = idx
		} else {
			// The row's activity isn't in the fetched list — make it selectable so
			// the entry lands on the intended activity, not a silent default.
			f.activities = append([]api.Activity{*prefActivity}, activities...)
			f.activityIdx = 0
		}
	}
	return f
}

// activityIndex returns the index of the activity with the given id, or -1 when
// it isn't present.
func activityIndex(activities []api.Activity, id uint64) int {
	for i, a := range activities {
		if a.ID == id {
			return i
		}
	}
	return -1
}

// cycleActivity advances the activity selector by delta, wrapping around.
func (f *addForm) cycleActivity(delta int) {
	n := len(f.activities)
	if n == 0 {
		return
	}
	f.activityIdx = ((f.activityIdx+delta)%n + n) % n
}

// focusedInput returns the text input for the focused field, or nil for the
// activity selector (which isn't a text field).
func (f *addForm) focusedInput() *textInput {
	switch f.field {
	case fieldTicket:
		return &f.ticket
	case fieldDuration:
		return &f.duration
	case fieldComment:
		return &f.comment
	default:
		return nil
	}
}

func (f addForm) ticketCode() string {
	return strings.TrimSpace(f.ticket.value)
}

// commentPtr returns the trimmed comment, or nil when empty so the backend stores
// a null comment rather than an empty string — matching `tl add`.
func (f addForm) commentPtr() *string {
	c := strings.TrimSpace(f.comment.value)
	if c == "" {
		return nil
	}
	return &c
}

// currentActivity returns the selected activity name for display, or "" when the
// list is empty.
func (f addForm) currentActivityName() string {
	if f.activityIdx < 0 || f.activityIdx >= len(f.activities) {
		return ""
	}
	return f.activities[f.activityIdx].Name
}
