package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

func newStatusCmd(r *root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the running timer, if any",
		Args:  cobra.NoArgs,
		// Pure local read: no auth or network, so status works offline and can
		// be used cheaply from a shell prompt.
		RunE: func(cmd *cobra.Command, _ []string) error {
			state, err := timer.Load()
			if errors.Is(err, timer.ErrNotRunning) {
				if r.json {
					fmt.Fprintln(cmd.OutOrStdout(), `{"running":false}`)
				} else {
					fmt.Fprintln(cmd.ErrOrStderr(), "no timer running")
				}
				return exitWith(1)
			}
			if err != nil {
				return err
			}

			elapsed := timer.ElapsedMinutes(state.Elapsed(time.Now()))

			if r.json {
				return printTimerJSON(cmd.OutOrStdout(), state, elapsed)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "⏱ %s (%s) — %s elapsed (started %s)\n",
				state.TicketCode, state.ActivityName, parse.FormatMinutes(elapsed),
				state.StartedAt.Format("15:04"))
			return nil
		},
	}
	return cmd
}

// printTimerJSON prints the state file plus the computed elapsed_minutes.
func printTimerJSON(w io.Writer, s timer.State, elapsedMinutes int) error {
	payload := struct {
		timer.State
		Running        bool `json:"running"`
		ElapsedMinutes int  `json:"elapsed_minutes"`
	}{State: s, Running: true, ElapsedMinutes: elapsedMinutes}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}
