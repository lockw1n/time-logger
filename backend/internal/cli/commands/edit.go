package commands

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

func newEditCmd(r *root) *cobra.Command {
	var (
		durationFlag string
		commentFlag  string
		yes          bool
	)

	cmd := &cobra.Command{
		Use:   "edit <entry-id>",
		Short: "Change the duration or comment of an entry",
		Long: strings.TrimSpace(`
Change the duration or comment of an existing entry (find its id with 'tl today').

Only the duration and comment can be edited. To change the ticket, activity or
date, delete the entry ('tl delete <id>') and add it again ('tl add').

Comments cannot be cleared through the API — the backend drops empty comments —
so to remove a comment, delete the entry and re-add it without one.`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseEntryID(args[0])
			if err != nil {
				return err
			}

			durSet := cmd.Flags().Changed("duration")
			comSet := cmd.Flags().Changed("comment")
			if !durSet && !comSet {
				fmt.Fprintln(cmd.ErrOrStderr(), "error: nothing to change — pass --duration and/or --comment")
				_ = cmd.Usage()
				return exitWith(1)
			}

			var req api.UpdateEntryRequest
			if durSet {
				minutes, err := parse.ParseDuration(durationFlag)
				if err != nil {
					return err
				}
				req.DurationMinutes = &minutes
			}
			if comSet {
				comment := strings.TrimSpace(commentFlag)
				if comment == "" {
					return errors.New("comments can't be cleared via the API — to remove a comment, 'tl delete' the entry and 'tl add' it again without one")
				}
				req.Comment = &comment
			}

			_, client, err := r.setup()
			if err != nil {
				return err
			}

			before, err := client.GetEntry(cmd.Context(), id)
			if err != nil {
				if notFoundErr(err) {
					return fmt.Errorf("entry %d not found", id)
				}
				return err
			}

			// The ticket lookup is only for the confirmation prompt, so skip its
			// extra request entirely when confirmation is skipped.
			if !yes {
				ticket := entryTicket(cmd.Context(), client, before)
				in := bufio.NewReader(cmd.InOrStdin())
				if !confirm(in, editSummary(id, ticket, before, req), false) {
					return errors.New("aborted")
				}
			}

			after, err := client.UpdateEntry(cmd.Context(), id, req)
			if err != nil {
				if notFoundErr(err) {
					return fmt.Errorf("entry %d not found", id)
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "updated entry #%d: %s\n", after.ID, entryOneLine(after))
			return nil
		},
	}

	cmd.Flags().StringVar(&durationFlag, "duration", "", "new duration (e.g. 2h, 90m, 1h30m)")
	cmd.Flags().StringVar(&commentFlag, "comment", "", "replace the comment")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")
	return cmd
}

// editSummary is the one-line before → after confirmation, listing only the
// fields the request actually changes.
func editSummary(id uint64, ticket string, before api.Entry, req api.UpdateEntryRequest) string {
	var parts []string
	if req.DurationMinutes != nil {
		parts = append(parts, fmt.Sprintf("%s → %s",
			parse.FormatMinutes(before.DurationMinutes), parse.FormatMinutes(*req.DurationMinutes)))
	}
	if req.Comment != nil {
		parts = append(parts, fmt.Sprintf("comment %s → %s",
			quoteComment(api.Deref(before.Comment)), quoteComment(*req.Comment)))
	}
	return fmt.Sprintf("edit #%d %s (%s): %s?", id, ticket, before.Date, strings.Join(parts, "; "))
}

// entryOneLine summarizes an entry's duration and comment for the success line.
func entryOneLine(e api.Entry) string {
	s := parse.FormatMinutes(e.DurationMinutes)
	if c := api.Deref(e.Comment); c != "" {
		s += " " + quoteComment(c)
	}
	return s
}
