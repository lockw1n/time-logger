package commands

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/auth"
	"github.com/lockw1n/time-logger/internal/cli/config"
)

func newLoginCmd(r *root) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate against the backend and store the session token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			in := bufio.NewReader(cmd.InOrStdin())

			email, err := promptEmail(in, cfg.Email)
			if err != nil {
				return err
			}

			password, err := promptPassword(in)
			if err != nil {
				return err
			}

			client, err := r.client(cfg)
			if err != nil {
				return err
			}

			result, err := client.Login(cmd.Context(), email, password)
			var apiErr *api.APIError
			if errors.As(err, &apiErr) && apiErr.Status == http.StatusUnauthorized {
				return errors.New("invalid credentials")
			}
			if err != nil {
				return err
			}

			if err := auth.Save(auth.Token{AccessToken: result.AccessToken, ExpiresAt: result.ExpiresAt}); err != nil {
				return err
			}

			if cfg.Email != email {
				cfg.Email = email
				if err := config.Save(cfg); err != nil {
					return err
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "logged in — session valid until %s\n", result.ExpiresAt.Local().Format("15:04"))
			return nil
		},
	}
}

func promptEmail(in *bufio.Reader, saved string) (string, error) {
	if saved != "" {
		fmt.Fprintf(os.Stderr, "email [%s]: ", saved)
	} else {
		fmt.Fprint(os.Stderr, "email: ")
	}

	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		return "", fmt.Errorf("reading email: %w", err)
	}

	email := strings.TrimSpace(line)
	if email == "" {
		email = saved
	}
	if email == "" {
		return "", errors.New("email is required")
	}

	return email, nil
}

func promptPassword(in *bufio.Reader) (string, error) {
	fmt.Fprint(os.Stderr, "password: ")

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		raw, err := term.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", fmt.Errorf("reading password: %w", err)
		}
		if len(raw) == 0 {
			return "", errors.New("password is required")
		}
		return string(raw), nil
	}

	// Non-interactive stdin (pipe, test): fall back to a plain line read.
	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		return "", fmt.Errorf("reading password: %w", err)
	}

	password := strings.TrimRight(line, "\r\n")
	if password == "" {
		return "", errors.New("password is required")
	}

	return password, nil
}
