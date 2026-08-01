package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/config"
)

func newPingCmd(r *root) *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Short: "Check that the backend is reachable (no auth required)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			client, err := r.client(cfg)
			if err != nil {
				return err
			}

			if err := client.Health(cmd.Context()); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "ok")
			return nil
		},
	}
}
