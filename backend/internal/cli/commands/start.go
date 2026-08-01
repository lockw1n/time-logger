package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/git"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

func newStartCmd(r *root) *cobra.Command {
	var (
		activityFlag string
		note         string
		atFlag       string
	)

	cmd := &cobra.Command{
		Use:   "start [ticket-code]",
		Short: "Start a live timer for a ticket",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now()

			// Fast path: surface an already-running timer — or an unreadable timer
			// file (whose error points at `tl cancel`) — before doing any network
			// work. The authoritative guard is the atomic create below.
			switch existing, err := timer.Load(); {
			case err == nil:
				return runningError(existing)
			case !errors.Is(err, timer.ErrNotRunning):
				return err
			}

			cfg, client, err := r.setup()
			if err != nil {
				return err
			}
			companyID, err := r.resolveCompany(cfg)
			if err != nil {
				return err
			}

			ticketCode, err := resolveTicketCode(cmd, args)
			if err != nil {
				return err
			}

			activityName, err := resolveActivityName(activityFlag, cfg)
			if err != nil {
				return err
			}
			// Resolve the activity now so `stop` can't fail hours later on a typo.
			activity, err := lookupActivity(cmd.Context(), client, companyID, activityName)
			if err != nil {
				return err
			}

			startedAt, err := resolveStartTime(atFlag, now)
			if err != nil {
				return err
			}

			state := timer.State{
				TicketCode:   ticketCode,
				ActivityID:   activity.ID,
				ActivityName: activity.Name,
				CompanyID:    companyID,
				StartedAt:    startedAt,
				Note:         strings.TrimSpace(note),
			}
			if err := timer.Start(state); err != nil {
				if errors.Is(err, timer.ErrRunning) {
					if existing, loadErr := timer.Load(); loadErr == nil {
						return runningError(existing)
					}
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "⏱ started %s (%s) at %s\n",
				ticketCode, activity.Name, startedAt.Format("15:04"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&activityFlag, "activity", "a", "", "activity name (required unless default_activity is set)")
	cmd.Flags().StringVar(&note, "note", "", "note attached to the entry when the timer stops")
	cmd.Flags().StringVar(&atFlag, "at", "", "backdate the start to today at HH:MM (must be in the past)")
	return cmd
}

// resolveTicketCode takes the explicit ticket argument, or infers it from the
// current git branch when none is given, printing what it inferred.
func resolveTicketCode(cmd *cobra.Command, args []string) (string, error) {
	if len(args) == 1 {
		code := strings.TrimSpace(args[0])
		if code == "" {
			return "", errors.New("ticket code is required")
		}
		return code, nil
	}

	branch, err := git.CurrentBranch(cmd.Context(), ".")
	if err != nil {
		return "", fmt.Errorf("%w — pass an explicit ticket code, e.g. 'tl start APP-123'", err)
	}
	code, ok := git.TicketFromBranch(branch)
	if !ok {
		return "", fmt.Errorf("no ticket code found in branch %q — pass one explicitly, e.g. 'tl start APP-123'", branch)
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "inferred ticket %s from branch '%s'\n", code, branch)
	return code, nil
}

// resolveStartTime returns now, or — when --at HH:MM is given — today at that
// time, which must be in the past (backdating a still-running timer forward
// would report negative elapsed time).
func resolveStartTime(atFlag string, now time.Time) (time.Time, error) {
	atFlag = strings.TrimSpace(atFlag)
	if atFlag == "" {
		return now, nil
	}

	t, err := time.ParseInLocation("15:04", atFlag, now.Location())
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --at time %q (expected HH:MM)", atFlag)
	}
	started := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	if !started.Before(now) {
		return time.Time{}, fmt.Errorf("--at %s is not in the past", atFlag)
	}
	return started, nil
}

func runningError(s timer.State) error {
	return fmt.Errorf("a timer for %s (%s) has been running since %s — stop or cancel it first",
		s.TicketCode, s.ActivityName, s.StartedAt.Format("15:04"))
}
