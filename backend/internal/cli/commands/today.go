package commands

import (
	"errors"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/render"
)

func newTodayCmd(r *root) *cobra.Command {
	var dateFlag string

	cmd := &cobra.Command{
		Use:   "today",
		Short: "Show the timesheet for a single day (default today)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, client, err := r.setup()
			if err != nil {
				return err
			}

			companyID, err := r.resolveCompany(cfg)
			if err != nil {
				return err
			}

			date, err := readDate(dateFlag)
			if err != nil {
				return err
			}

			if r.json {
				raw, err := client.GetRaw(cmd.Context(), api.TimesheetPath(companyID, date, date))
				if err != nil {
					return err
				}
				return printRawJSON(cmd.OutOrStdout(), raw)
			}

			ts, err := client.Timesheet(cmd.Context(), companyID, date, date)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(render.Today(ts)))
			return err
		},
	}

	cmd.Flags().StringVarP(&dateFlag, "date", "d", "", "day to show (YYYY-MM-DD, today, yesterday, mon-sun)")
	return cmd
}

// readDate normalizes a date flag for read commands. Unlike `add`, viewing a
// future day is harmless, so ErrFutureDate is not treated as a problem.
func readDate(flag string) (string, error) {
	date, err := parse.ParseDate(flag, time.Now())
	if err != nil && !errors.Is(err, parse.ErrFutureDate) {
		return "", err
	}
	return date, nil
}
