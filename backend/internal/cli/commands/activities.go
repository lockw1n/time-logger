package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

func newActivitiesCmd(r *root) *cobra.Command {
	return &cobra.Command{
		Use:   "activities",
		Short: "List the activities available for the resolved company",
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

			if r.json {
				raw, err := client.GetRaw(cmd.Context(), api.ActivitiesPath(companyID))
				if err != nil {
					return err
				}
				return printRawJSON(cmd.OutOrStdout(), raw)
			}

			// Use the cache-warming lookup so a successful online run populates the
			// activities cache that offline `tl add` depends on (and so this
			// command itself still works from cache when offline).
			activities, err := activitiesForLookup(cmd.Context(), client, companyID)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "id\tname\tbillable\tpriority")
			for _, a := range activities {
				fmt.Fprintf(w, "%d\t%s\t%t\t%d\n", a.ID, a.Name, a.Billable, a.Priority)
			}
			return w.Flush()
		},
	}
}
