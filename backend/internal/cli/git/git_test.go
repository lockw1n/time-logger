package git

import "testing"

func TestTicketFromBranch(t *testing.T) {
	cases := []struct {
		branch string
		want   string
		ok     bool
	}{
		{"feature/APP-123-fix-mapper", "APP-123", true},
		{"APP-123", "APP-123", true},
		{"bugfix/app-99-x", "APP-99", true},
		{"feature/proj-7", "PROJ-7", true},
		{"main", "", false},
		{"release/2.0", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := TicketFromBranch(c.branch)
		if got != c.want || ok != c.ok {
			t.Errorf("TicketFromBranch(%q) = (%q, %v), want (%q, %v)", c.branch, got, ok, c.want, c.ok)
		}
	}
}
