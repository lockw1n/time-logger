package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/outbox"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

func newSyncCmd(r *root) *cobra.Command {
	var (
		list    bool
		discard string
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Replay queued offline writes to the backend",
		Long: "Replay time entries that were queued while the backend was unreachable.\n\n" +
			"Running any other command already syncs in the background; use 'tl sync'\n" +
			"to force a foreground run, inspect the queue (--list), or drop a permanently\n" +
			"failed op (--discard <file>).",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch {
			case list:
				return runSyncList(cmd)
			case discard != "":
				return runSyncDiscard(cmd, discard)
			default:
				return runSyncReplay(cmd, r)
			}
		},
	}

	cmd.Flags().BoolVar(&list, "list", false, "show pending and failed ops without replaying")
	cmd.Flags().StringVar(&discard, "discard", "", "remove a failed op by filename (see --list)")
	return cmd
}

// runSyncReplay runs the outbox replay in the foreground and reports the result,
// exiting non-zero when anything is left pending or was moved to failed/.
func runSyncReplay(cmd *cobra.Command, r *root) error {
	_, client, err := r.loadClient()
	if err != nil {
		return err
	}

	rep, err := outbox.Sync(cmd.Context(), client, outbox.Options{Notice: cmd.ErrOrStderr()})
	if errors.Is(err, outbox.ErrLocked) {
		fmt.Fprintln(cmd.ErrOrStderr(), "another sync is already running")
		return exitWith(1)
	}
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if rep.Synced == 0 && rep.Failed == 0 && rep.Remaining == 0 {
		fmt.Fprintln(out, "outbox is empty — nothing to sync")
		return nil
	}
	fmt.Fprintf(out, "synced %d, failed %d, remaining %d\n", rep.Synced, rep.Failed, rep.Remaining)

	if rep.Remaining > 0 || rep.Failed > 0 {
		return exitWith(1)
	}
	return nil
}

// runSyncList prints the pending and failed queues without replaying anything.
func runSyncList(cmd *cobra.Command) error {
	pending, err := outbox.Pending()
	if err != nil {
		return err
	}
	failed, err := outbox.Failed()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if len(pending) == 0 && len(failed) == 0 {
		fmt.Fprintln(out, "outbox is empty")
		return nil
	}

	if len(pending) > 0 {
		fmt.Fprintf(out, "pending (%d):\n", len(pending))
		writeOpTable(out, pending, false)
	}
	if len(failed) > 0 {
		if len(pending) > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "failed (%d):\n", len(failed))
		writeOpTable(out, failed, true)
	}
	return nil
}

// writeOpTable renders queued ops as an aligned table. The last_error column is
// only shown for the failed queue, where it explains why an op is stuck.
func writeOpTable(out io.Writer, entries []outbox.Entry, withError bool) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	header := "  FILE\tTICKET\tDATE\tMINUTES\tATTEMPTS"
	if withError {
		header += "\tLAST ERROR"
	}
	fmt.Fprintln(tw, header)
	for _, e := range entries {
		row := fmt.Sprintf("  %s\t%s\t%s\t%s\t%d",
			e.Filename, e.Op.Payload.TicketCode, e.Op.Payload.Date,
			parse.FormatMinutes(e.Op.Payload.DurationMinutes), e.Op.Attempts)
		if withError {
			row += "\t" + e.Op.LastError
		}
		fmt.Fprintln(tw, row)
	}
	_ = tw.Flush()
}

// runSyncDiscard removes a permanently-failed op after confirmation.
func runSyncDiscard(cmd *cobra.Command, filename string) error {
	in := bufio.NewReader(cmd.InOrStdin())
	if !confirm(in, fmt.Sprintf("permanently discard failed op %s?", filename), false) {
		return errors.New("aborted")
	}
	if err := outbox.Discard(filename); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "discarded %s\n", filename)
	return nil
}
