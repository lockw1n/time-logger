package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

func newInvoiceCmd(r *root) *cobra.Command {
	var (
		monthFlag  string
		startFlag  string
		endFlag    string
		formatFlag string
		outputFlag string
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "Generate an invoice document for a period and write it to disk",
		Long: strings.TrimSpace(`
Generate an invoice (PDF or Excel) for a period and write it to disk.

Give the period either as a whole month (--month YYYY-MM) or an explicit range
(--start / --end); the two forms are mutually exclusive. With neither, the
previous calendar month is used, since invoices are generated after a period closes.`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			start, end, err := resolveInvoicePeriod(monthFlag, startFlag, endFlag, time.Now())
			if err != nil {
				return err
			}

			format := strings.ToLower(strings.TrimSpace(formatFlag))
			if format != "pdf" && format != "excel" {
				return fmt.Errorf("invalid --format %q (accepted: pdf, excel)", formatFlag)
			}

			cfg, client, err := r.setup()
			if err != nil {
				return err
			}
			companyID, err := r.resolveCompany(cfg)
			if err != nil {
				return err
			}

			// The PDF render round-trip is slow; explain the wait before blocking.
			fmt.Fprintln(cmd.ErrOrStderr(), "generating…")

			file, err := client.GenerateInvoice(cmd.Context(), api.GenerateInvoiceRequest{
				CompanyID: companyID,
				Start:     start,
				End:       end,
				Format:    format,
			})
			if err != nil {
				return err
			}

			path, err := invoiceOutputPath(outputFlag, file.Filename)
			if err != nil {
				return err
			}
			if !force {
				if _, err := os.Stat(path); err == nil {
					return fmt.Errorf("%s already exists — pass --force to overwrite", path)
				}
			}

			if err := os.WriteFile(path, file.Data, 0o644); err != nil {
				return fmt.Errorf("writing invoice: %w", err)
			}

			abs, err := filepath.Abs(path)
			if err != nil {
				abs = path
			}
			fmt.Fprintln(cmd.OutOrStdout(), abs)
			return nil
		},
	}

	cmd.Flags().StringVar(&monthFlag, "month", "", "period as a whole month (YYYY-MM); default: previous month")
	cmd.Flags().StringVar(&startFlag, "start", "", "period start (YYYY-MM-DD); requires --end")
	cmd.Flags().StringVar(&endFlag, "end", "", "period end (YYYY-MM-DD); requires --start")
	cmd.Flags().StringVar(&formatFlag, "format", "pdf", "document format: pdf or excel")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "output file or directory (default: current directory)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}

// resolveInvoicePeriod turns the --month / --start / --end flags into a concrete
// YYYY-MM-DD range. The two forms are mutually exclusive; --start and --end must
// be given together; with neither form the previous calendar month is used.
func resolveInvoicePeriod(month, start, end string, now time.Time) (string, string, error) {
	month = strings.TrimSpace(month)
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)

	if month != "" && (start != "" || end != "") {
		return "", "", errors.New("--month is mutually exclusive with --start/--end")
	}
	if (start == "") != (end == "") {
		return "", "", errors.New("--start and --end must be given together")
	}

	if start != "" {
		if err := validateExplicitRange(start, end); err != nil {
			return "", "", err
		}
		return start, end, nil
	}

	if month == "" {
		month = parse.PreviousMonth(now)
	}
	return parse.MonthRange(month)
}

// validateExplicitRange rejects malformed or reversed --start/--end before the
// slow server round-trip, so a typo fails fast with a clear message instead of an
// opaque backend error after the "generating…" wait.
func validateExplicitRange(start, end string) error {
	s, err := time.Parse(parse.DateLayout, start)
	if err != nil {
		return fmt.Errorf("invalid --start %q (expected YYYY-MM-DD)", start)
	}
	e, err := time.Parse(parse.DateLayout, end)
	if err != nil {
		return fmt.Errorf("invalid --end %q (expected YYYY-MM-DD)", end)
	}
	if e.Before(s) {
		return fmt.Errorf("--end %s is before --start %s", end, start)
	}
	return nil
}

// invoiceOutputPath resolves where to write the document. With no -o the server
// filename lands in the current directory; -o pointing at an existing directory
// (or a trailing-slash path) puts the server filename inside it; any other -o is
// treated as the full target path.
func invoiceOutputPath(output, filename string) (string, error) {
	if output == "" {
		return filename, nil
	}
	if info, err := os.Stat(output); err == nil && info.IsDir() {
		return filepath.Join(output, filename), nil
	}
	if strings.HasSuffix(output, string(os.PathSeparator)) {
		return filepath.Join(output, filename), nil
	}
	return output, nil
}
