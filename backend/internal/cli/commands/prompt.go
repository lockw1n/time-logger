package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// confirm asks a yes/no question on stderr and reads the answer from in.
// defaultYes controls both the shown default ([Y/n] vs [y/N]) and the response
// to an empty line (Enter). A closed/non-interactive stdin (EOF with nothing
// typed) is NOT treated as accepting the default: it returns false so a mutating
// prompt never proceeds unattended — callers must pass --yes to skip it.
func confirm(in *bufio.Reader, prompt string, defaultYes bool) bool {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}
	fmt.Fprintf(os.Stderr, "%s %s ", prompt, suffix)

	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		return false
	}

	switch strings.ToLower(strings.TrimSpace(line)) {
	case "":
		return defaultYes
	case "y", "yes":
		return true
	default:
		return false
	}
}
