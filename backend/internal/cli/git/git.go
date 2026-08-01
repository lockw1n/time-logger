// Package git provides the minimal git integration the CLI needs: reading the
// current branch and inferring a JIRA-style ticket code from it, so `tl start`
// on a feature branch can default the ticket without the user retyping it.
package git

import (
	"context"
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

// ticketRE matches a ticket code anywhere in a branch name, case-insensitively:
// a letter, one or more letters/digits, a hyphen, then digits (e.g. APP-123).
var ticketRE = regexp.MustCompile(`(?i)[a-z][a-z0-9]+-[0-9]+`)

// ErrDetached is returned by CurrentBranch when HEAD is detached (no branch name).
var ErrDetached = errors.New("detached HEAD")

// CurrentBranch returns the abbreviated current branch in dir, or an error when
// dir is not a git work tree, git is unavailable, or HEAD is detached.
func CurrentBranch(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("not a git repository (or git is unavailable)")
	}
	branch := strings.TrimSpace(string(out))
	if branch == "" || branch == "HEAD" {
		return "", ErrDetached
	}
	return branch, nil
}

// TicketFromBranch extracts the first ticket code from a branch name and
// uppercases it (branches are commonly lowercase). ok is false when the branch
// carries no recognizable code.
func TicketFromBranch(branch string) (string, bool) {
	m := ticketRE.FindString(branch)
	if m == "" {
		return "", false
	}
	return strings.ToUpper(m), true
}
