package commands

import (
	"bufio"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

func newCancelCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Discard the running timer without logging it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			state, err := timer.Load()
			if errors.Is(err, timer.ErrNotRunning) {
				fmt.Fprintln(cmd.ErrOrStderr(), "no timer running")
				return exitWith(1)
			}

			in := bufio.NewReader(cmd.InOrStdin())

			// A corrupt timer file can't be described, but cancel exists precisely
			// to clear it — offer to delete it directly.
			var corrupt *timer.CorruptError
			if errors.As(err, &corrupt) {
				if !yes && !confirm(in, "timer file is unreadable — discard it?", false) {
					return errors.New("aborted")
				}
				return timer.Clear()
			}
			if err != nil {
				return err
			}

			elapsed := timer.ElapsedMinutes(state.Elapsed(time.Now()))
			if !yes && !confirm(in, fmt.Sprintf("discard %s on %s?", parse.FormatMinutes(elapsed), state.TicketCode), false) {
				return errors.New("aborted")
			}

			if err := timer.Clear(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "discarded %s on %s\n", parse.FormatMinutes(elapsed), state.TicketCode)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")
	return cmd
}
