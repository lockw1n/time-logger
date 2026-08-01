package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/auth"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Delete the stored session token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := auth.Delete(); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "logged out")
			return nil
		},
	}
}
