package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

func newCompaniesCmd(r *root) *cobra.Command {
	return &cobra.Command{
		Use:   "companies",
		Short: "List the companies you can log time to",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, client, err := r.setup()
			if err != nil {
				return err
			}

			if r.json {
				raw, err := client.GetRaw(cmd.Context(), api.CompaniesPath())
				if err != nil {
					return err
				}
				return printRawJSON(cmd.OutOrStdout(), raw)
			}

			companies, err := client.Companies(cmd.Context())
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "id\tname\tname_short")
			for _, c := range companies {
				fmt.Fprintf(w, "%d\t%s\t%s\n", c.ID, c.Name, api.Deref(c.NameShort))
			}
			return w.Flush()
		},
	}
}
