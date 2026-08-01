package commands

import (
	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/tui"
)

// newUICmd wires up `tl ui`, the full-screen weekly timesheet. It resolves the
// company and opening week the same way the one-shot read commands do, then hands
// off to the Bubble Tea program, which owns the terminal until the user quits.
func newUICmd(r *root) *cobra.Command {
	var (
		dateFlag  string
		panicTest bool
	)

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Open the interactive weekly timesheet",
		Long: "Open a full-screen, editable weekly timesheet grid.\n\n" +
			"Navigate cells with the arrows (or hjkl); enter edits a cell, a adds an\n" +
			"entry, d deletes one; [ and ] page weeks, t jumps to this week, c cycles\n" +
			"company, r refreshes. q or Ctrl-C quits.",
		Args: cobra.NoArgs,
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

			return tui.Run(cmd.Context(), client, tui.Options{
				CompanyID:  companyID,
				Date:       date,
				PanicOnKey: panicTest,
			})
		},
	}

	cmd.Flags().StringVarP(&dateFlag, "date", "d", "", "any day in the week to open (YYYY-MM-DD, today, yesterday, mon-sun)")
	cmd.Flags().BoolVar(&panicTest, "panic-test", false, "panic on the next keypress to verify terminal restoration (manual testing)")
	_ = cmd.Flags().MarkHidden("panic-test")
	return cmd
}
