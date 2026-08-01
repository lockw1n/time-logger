package commands

import (
	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
	"github.com/lockw1n/time-logger/internal/cli/render"
)

func newWeekCmd(r *root) *cobra.Command {
	var dateFlag string

	cmd := &cobra.Command{
		Use:   "week",
		Short: "Show the Monday–Sunday timesheet grid for a week (default this week)",
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

			monday, sunday, err := parse.WeekBounds(date)
			if err != nil {
				return err
			}

			if r.json {
				raw, err := client.GetRaw(cmd.Context(), api.TimesheetPath(companyID, monday, sunday))
				if err != nil {
					return err
				}
				return printRawJSON(cmd.OutOrStdout(), raw)
			}

			ts, err := client.Timesheet(cmd.Context(), companyID, monday, sunday)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(render.Week(ts)))
			return err
		},
	}

	cmd.Flags().StringVarP(&dateFlag, "date", "d", "", "any day in the target week (YYYY-MM-DD, today, yesterday, mon-sun)")
	return cmd
}
