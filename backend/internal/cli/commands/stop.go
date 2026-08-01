package commands

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/auth"
	"github.com/lockw1n/time-logger/internal/cli/entrylog"
	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

// twelveHours (in minutes) is the point past which a *measured* run is treated
// as a probably-forgotten timer and confirmed before logging.
const twelveHours = 12 * 60

func newStopCmd(r *root) *cobra.Command {
	var (
		durationFlag string
		yes          bool
	)

	cmd := &cobra.Command{
		Use:   "stop [comment...]",
		Short: "Stop the running timer and log the elapsed time",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := timer.Load()
			if errors.Is(err, timer.ErrNotRunning) {
				fmt.Fprintln(cmd.ErrOrStderr(), "no timer running")
				return exitWith(1)
			}
			if err != nil {
				return err
			}

			now := time.Now()
			measured := timer.ElapsedMinutes(state.Elapsed(now))

			minutes := measured
			overridden := false
			if strings.TrimSpace(durationFlag) != "" {
				minutes, err = parse.ParseDuration(durationFlag)
				if err != nil {
					return err
				}
				overridden = true
			}

			// A measured run past the 24h/entry ceiling can never be logged as-is
			// (the backend rejects it), so fail early with the way out instead of
			// letting the server bounce it — --duration is the intended escape.
			if !overridden && measured > parse.MaxMinutes {
				return fmt.Errorf(
					"timer ran %s, over the 24h per-entry limit — pass --duration to log a specific amount (the timer is kept)",
					parse.FormatMinutes(measured))
			}

			in := bufio.NewReader(cmd.InOrStdin())

			// Guard forgotten timers: a very long measured run is almost always a
			// timer left running overnight, so confirm before logging it.
			if !overridden && measured > twelveHours && !yes {
				if !confirm(in, fmt.Sprintf("timer ran %s — really log it?", parse.FormatMinutes(measured)), false) {
					return errors.New("aborted — timer left running")
				}
			}

			_, client, err := r.setup()
			if err != nil {
				return err
			}

			// A timer always logs to the day it is stopped.
			date := now.Format(parse.DateLayout)

			// Use the activity id captured at start so no network lookup is needed
			// here — that is what lets stop queue the entry when offline. Timers
			// started before ids were stored fall back to a name lookup (online).
			activity := api.Activity{ID: state.ActivityID, Name: state.ActivityName}
			if activity.ID == 0 {
				activity, err = lookupActivity(cmd.Context(), client, state.CompanyID, state.ActivityName)
				if err != nil {
					return stopFailed(cmd, err)
				}
			}

			// Comment precedence: the stop argument wins over the start --note.
			comment := commentPtr(args)
			if comment == nil && state.Note != "" {
				n := state.Note
				comment = &n
			}

			err = submitEntry(cmd.Context(), client, entrylog.Request{
				CompanyID:  state.CompanyID,
				TicketCode: state.TicketCode,
				Activity:   activity,
				Date:       date,
				Minutes:    minutes,
				Comment:    comment,
			}, entrylog.Options{
				Out:         cmd.OutOrStdout(),
				Confirm:     func(p string, d bool) bool { return confirm(in, p, d) },
				SkipConfirm: yes,
			})
			// A queued create captured the measured time in the payload, so the
			// timer has done its job — clear it and propagate the exit-3 signal.
			// Any other failure keeps the timer so nothing is lost on retry.
			if err != nil && !isQueued(err) {
				return stopFailed(cmd, err)
			}
			if clearErr := timer.Clear(); clearErr != nil {
				// The entry is already logged (or queued); a failure to remove the
				// timer file must not masquerade as a submit failure — that would
				// invite a retry that double-logs the same span. Warn, point at the
				// manual fix, and keep the real outcome (nil, or exit-3 queued).
				fmt.Fprintf(cmd.ErrOrStderr(),
					"warning: entry recorded but the timer could not be cleared (%v) — run 'tl cancel'\n", clearErr)
			}
			return err
		},
	}

	cmd.Flags().StringVar(&durationFlag, "duration", "", "log this duration instead of the measured elapsed time")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmations")
	return cmd
}

// stopFailed handles any failure after we've committed to logging. The timer is
// never cleared on this path, so nothing is lost. Auth failures get the
// login-and-retry message and exit 2; every other failure keeps the timer with a
// generic retry hint and exits 1.
func stopFailed(cmd *cobra.Command, err error) error {
	if errors.Is(err, auth.ErrNotLoggedIn) || errors.Is(err, api.ErrUnauthorized) {
		fmt.Fprintln(cmd.ErrOrStderr(), "timer kept — run 'tl login' then 'tl stop' again")
		return exitWith(2)
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\ntimer kept — fix the problem and run 'tl stop' again\n", err)
	return exitWith(1)
}
