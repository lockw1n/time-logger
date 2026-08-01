package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/auth"
	"github.com/lockw1n/time-logger/internal/cli/commands"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	err := commands.Execute(version)
	if err == nil {
		return
	}

	// A command may take full control of its exit code and its own output.
	var coder interface{ ExitCode() int }
	if errors.As(err, &coder) {
		os.Exit(coder.ExitCode())
	}

	if errors.Is(err, auth.ErrNotLoggedIn) || errors.Is(err, api.ErrUnauthorized) {
		fmt.Fprintln(os.Stderr, "error: not logged in or session expired — run 'tl login'")
		os.Exit(2)
	}

	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}
