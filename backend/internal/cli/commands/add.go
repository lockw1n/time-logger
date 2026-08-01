package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/config"
	"github.com/lockw1n/time-logger/internal/cli/entrylog"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

func newAddCmd(r *root) *cobra.Command {
	var (
		dateFlag     string
		activityFlag string
		yes          bool
	)

	cmd := &cobra.Command{
		Use:   "add <ticket-code> <duration> [comment...]",
		Short: "Log time to a ticket",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := r.setup()
			if err != nil {
				return err
			}

			companyID, err := r.resolveCompany(cfg)
			if err != nil {
				return err
			}

			ticketCode := strings.TrimSpace(args[0])
			if ticketCode == "" {
				return errors.New("ticket code is required")
			}

			minutes, err := parse.ParseDuration(args[1])
			if err != nil {
				return err
			}

			comment := commentPtr(args[2:])

			in := bufio.NewReader(cmd.InOrStdin())

			date, err := resolveAddDate(dateFlag, yes, in)
			if err != nil {
				return err
			}

			activityName, err := resolveActivityName(activityFlag, cfg)
			if err != nil {
				return err
			}

			activity, err := lookupActivity(cmd.Context(), client, companyID, activityName)
			if err != nil {
				// Offline logging needs the activity id, which comes from a cached
				// activities list. With no cache we can't build a payload to queue,
				// so explain the prerequisite instead of leaking a bare dial error.
				if errors.Is(err, api.ErrUnreachable) {
					return fmt.Errorf(
						"backend unreachable and no cached activities for company %d — run 'tl activities' once while online so offline logging can resolve %q",
						companyID, activityName)
				}
				return err
			}

			return submitEntry(cmd.Context(), client, entrylog.Request{
				CompanyID:  companyID,
				TicketCode: ticketCode,
				Activity:   activity,
				Date:       date,
				Minutes:    minutes,
				Comment:    comment,
			}, entrylog.Options{
				Out:         cmd.OutOrStdout(),
				Confirm:     func(p string, d bool) bool { return confirm(in, p, d) },
				SkipConfirm: yes,
			})
		},
	}

	cmd.Flags().StringVarP(&dateFlag, "date", "d", "", "day to log (YYYY-MM-DD, today, yesterday, mon-sun)")
	cmd.Flags().StringVarP(&activityFlag, "activity", "a", "", "activity name (required unless default_activity is set)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmations")
	return cmd
}

// commentPtr joins the trailing comment args, returning nil when empty so the
// backend stores a null comment rather than an empty string.
func commentPtr(args []string) *string {
	c := strings.TrimSpace(strings.Join(args, " "))
	if c == "" {
		return nil
	}
	return &c
}

// resolveAddDate parses the date flag and, for a future date, confirms unless
// --yes was given (logging future time is almost always a typo).
func resolveAddDate(flag string, skipConfirm bool, in *bufio.Reader) (string, error) {
	date, err := parse.ParseDate(flag, time.Now())
	if errors.Is(err, parse.ErrFutureDate) {
		if skipConfirm {
			return date, nil
		}
		if !confirm(in, fmt.Sprintf("%s is in the future — log time there anyway?", date), false) {
			return "", errors.New("aborted")
		}
		return date, nil
	}
	if err != nil {
		return "", err
	}
	return date, nil
}

// resolveActivityName picks the activity name from the flag, falling back to
// default_activity in config.
func resolveActivityName(flag string, cfg config.Config) (string, error) {
	name := strings.TrimSpace(flag)
	if name == "" {
		name = strings.TrimSpace(cfg.DefaultActivity)
	}
	if name == "" {
		return "", errors.New("no activity — pass --activity <name> or run 'tl config set default_activity <name>' (see 'tl activities')")
	}
	return name, nil
}

// lookupActivity resolves an activity name to its record for the company,
// matching case-insensitively. On a miss or ambiguity it lists the available
// activity names so the user can correct the input.
func lookupActivity(ctx context.Context, client api.Client, companyID uint64, name string) (api.Activity, error) {
	activities, err := activitiesForLookup(ctx, client, companyID)
	if err != nil {
		return api.Activity{}, err
	}

	var matches []api.Activity
	for _, a := range activities {
		if strings.EqualFold(a.Name, name) {
			matches = append(matches, a)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return api.Activity{}, fmt.Errorf("unknown activity %q — available: %s", name, activityNames(activities))
	default:
		return api.Activity{}, fmt.Errorf("ambiguous activity %q matches multiple activities: %s", name, activityNames(matches))
	}
}

func activityNames(activities []api.Activity) string {
	names := make([]string, 0, len(activities))
	for _, a := range activities {
		names = append(names, a.Name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		return "(none)"
	}
	return strings.Join(names, ", ")
}
