package commands

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

func newDeleteCmd(r *root) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <entry-id>",
		Short: "Delete an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseEntryID(args[0])
			if err != nil {
				return err
			}

			_, client, err := r.setup()
			if err != nil {
				return err
			}

			entry, err := client.GetEntry(cmd.Context(), id)
			if err != nil {
				if notFoundErr(err) {
					return fmt.Errorf("entry %d not found", id)
				}
				return err
			}

			// The ticket lookup only feeds the confirmation prompt, so skip its
			// extra request entirely when confirmation is skipped.
			if !yes {
				ticket := entryTicket(cmd.Context(), client, entry)
				prompt := fmt.Sprintf("delete %s on %s (%s)%s?",
					parse.FormatMinutes(entry.DurationMinutes), ticket, entry.Date, commentSuffix(entry.Comment))
				in := bufio.NewReader(cmd.InOrStdin())
				if !confirm(in, prompt, false) {
					return errors.New("aborted")
				}
			}

			if err := client.DeleteEntry(cmd.Context(), id); err != nil {
				if notFoundErr(err) {
					return fmt.Errorf("entry %d not found", id)
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted entry #%d\n", id)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")
	return cmd
}

// commentSuffix renders a trailing ` "comment"` for the delete prompt, or "" when
// the entry has no comment.
func commentSuffix(c *string) string {
	if s := api.Deref(c); s != "" {
		return " " + quoteComment(s)
	}
	return ""
}
