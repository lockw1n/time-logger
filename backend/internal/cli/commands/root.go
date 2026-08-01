package commands

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/auth"
	"github.com/lockw1n/time-logger/internal/cli/config"
	"github.com/lockw1n/time-logger/internal/cli/outbox"
)

const defaultBaseURL = "http://localhost:8085"

// syncBudget bounds how long a command waits, before exiting, for a background
// replay to finish. Past it the sync is cancelled and reported as incomplete —
// a slow sync must never hold up the command the user actually ran.
const syncBudget = 5 * time.Second

// cancelGrace bounds the extra wait after the budget expires and the sync is
// cancelled, giving an in-flight op a brief chance to unwind cleanly. If it
// overruns even this, we report from the on-disk queue and let the goroutine
// finish detached, so the command's exit is bounded by syncBudget+cancelGrace.
const cancelGrace = 1 * time.Second

// root holds the values of the persistent flags shared by all subcommands, plus
// the handle to any background outbox sync started for this invocation.
type root struct {
	baseURL  string
	company  uint64
	json     bool
	caCert   string
	insecure bool

	// ctx and errOut are captured in the persistent pre-run so setup() (which
	// has neither) can derive the sync context and report to the right stream.
	ctx    context.Context
	errOut io.Writer

	syncHandle *backgroundSync
	waitOnce   sync.Once
}

// backgroundSync tracks the goroutine replaying the outbox alongside the
// foreground command. Its report, err and notice are written only by that
// goroutine and read only after done is closed, so no locking is needed.
type backgroundSync struct {
	cancel context.CancelFunc
	done   chan struct{}
	report outbox.Report // read only after done is closed
	err    error         // read only after done is closed
	// notice buffers the sync's own messages (e.g. breaking a stale lock) so the
	// background goroutine never writes the same stream as the foreground command;
	// it is flushed to stderr after the goroutine finishes.
	notice bytes.Buffer
}

// NewRoot builds the tl root command. Callers that also need the background-sync
// teardown (production, via Execute) use buildRoot to reach the shared state.
func NewRoot(version string) *cobra.Command {
	cmd, _ := buildRoot(version)
	return cmd
}

func buildRoot(version string) (*cobra.Command, *root) {
	r := &root{}

	cmd := &cobra.Command{
		Use:           "tl",
		Short:         "Terminal client for the time-logger backend",
		SilenceUsage:  true,
		SilenceErrors: true,
		// Capture the command's context and error stream so setup() can launch
		// the background sync against them. This runs for every command but does
		// no work itself — only API commands (via setup) actually start a sync.
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			r.ctx = cmd.Context()
			r.errOut = cmd.ErrOrStderr()
		},
		// After a command's own work, wait briefly for the background sync and
		// report it. This runs only when the command succeeded; the Execute
		// wrapper waits on the error path too, and waitForSync is idempotent.
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			r.waitForSync()
		},
	}

	cmd.PersistentFlags().StringVar(&r.baseURL, "base-url", "", "backend base URL (default: base_url from config, then "+defaultBaseURL+")")
	cmd.PersistentFlags().Uint64VarP(&r.company, "company", "c", 0, "company ID (overrides default_company from config)")
	cmd.PersistentFlags().BoolVar(&r.json, "json", false, "print raw API responses")
	cmd.PersistentFlags().StringVar(&r.caCert, "ca-cert", "", "PEM cert to trust for an HTTPS backend (overrides ca_cert from config)")
	cmd.PersistentFlags().BoolVar(&r.insecure, "insecure", false, "skip TLS verification for an HTTPS backend")

	cmd.AddCommand(
		newVersionCmd(version),
		newLoginCmd(r),
		newLogoutCmd(),
		newPingCmd(r),
		newConfigCmd(),
		newAddCmd(r),
		newTodayCmd(r),
		newWeekCmd(r),
		newCompaniesCmd(r),
		newActivitiesCmd(r),
		newStartCmd(r),
		newStatusCmd(r),
		newStopCmd(r),
		newCancelCmd(),
		newEditCmd(r),
		newDeleteCmd(r),
		newInvoiceCmd(r),
		newSyncCmd(r),
		newUICmd(r),
	)

	return cmd, r
}

// Execute wires signal handling and background-sync teardown around the root
// command, then runs it. SIGINT/SIGTERM cancels the shared context, which stops
// both the foreground command and any in-flight background sync; an unconfirmed
// op stays queued (at-least-once). main delegates to this rather than Execute()
// so the sync goroutine is always waited on, even when the command errored.
func Execute(version string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd, r := buildRoot(version)
	err := cmd.ExecuteContext(ctx)
	// Idempotent with the PersistentPostRun: covers the error path, where cobra
	// skips PostRun but a background sync may still be running.
	r.waitForSync()
	return err
}

// setup loads config and builds the API client — the boilerplate every
// backend-touching command starts with — and opportunistically launches a
// background replay of any queued offline writes.
func (r *root) setup() (config.Config, api.Client, error) {
	cfg, client, err := r.loadClient()
	if err != nil {
		return config.Config{}, nil, err
	}
	r.startBackgroundSync(client)
	return cfg, client, nil
}

// loadClient loads config and builds the API client without touching the outbox.
// `tl sync` uses this directly so a manual foreground sync doesn't also spawn a
// competing background one.
func (r *root) loadClient() (config.Config, api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return config.Config{}, nil, err
	}
	client, err := r.client(cfg)
	if err != nil {
		return config.Config{}, nil, err
	}
	return cfg, client, nil
}

// startBackgroundSync launches an outbox replay in a goroutine when there is
// work to do and a usable session, letting the foreground command proceed
// concurrently (the DB serializes the racing writes server-side). It is a no-op
// when the queue is empty or the user isn't logged in, so the common path pays
// only a cheap directory count.
func (r *root) startBackgroundSync(client api.Client) {
	if r.ctx == nil || r.syncHandle != nil {
		return
	}
	if n, err := outbox.Count(); err != nil || n == 0 {
		return
	}
	// Replaying needs a valid token; if the session is gone, leave the queue for
	// after the user logs in rather than burning an attempt on a guaranteed 401.
	if _, err := auth.Load(); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(r.ctx)
	h := &backgroundSync{cancel: cancel, done: make(chan struct{})}
	r.syncHandle = h
	go func() {
		defer close(h.done)
		h.report, h.err = outbox.Sync(ctx, client, outbox.Options{Notice: &h.notice})
	}()
}

// waitForSync blocks (once) for the background sync to finish, up to syncBudget,
// then reports the outcome to stderr. On timeout it cancels the sync and reports
// it as incomplete rather than making the user wait on a slow replay.
func (r *root) waitForSync() {
	r.waitOnce.Do(func() {
		h := r.syncHandle
		if h == nil {
			return
		}
		// Always release the derived context, including the fast happy path below
		// where the goroutine finishes before the budget — otherwise the
		// context's cancel func is never called and it leaks. Cancelling an
		// already-finished sync is a no-op, so this is safe on every path.
		defer h.cancel()
		select {
		case <-h.done:
			// Finished within budget: safe to read h.* (done closed → happens-before).
			r.reportSync(h.report, h.err, h.notice.String())
			return
		case <-time.After(syncBudget):
		}

		// Budget expired: ask the sync to stop and give it a brief grace period to
		// unwind. If it finishes, report normally; if it overruns even the grace,
		// report from the on-disk count WITHOUT touching h.* (still being written)
		// and let the goroutine finish detached — the command must not block on it.
		h.cancel()
		select {
		case <-h.done:
			r.reportSync(h.report, h.err, h.notice.String())
		case <-time.After(cancelGrace):
			n, _ := outbox.Count()
			fmt.Fprintf(r.stderr(), "sync incomplete — %s still queued\n", pluralOps(n))
		}
	})
}

// stderr returns the writer sync reports go to, defaulting to os.Stderr when the
// pre-run hasn't captured the command's error stream.
func (r *root) stderr() io.Writer {
	if r.errOut != nil {
		return r.errOut
	}
	return os.Stderr
}

// reportSync prints a one-line summary of a completed (or cancelled) replay.
// Silence when nothing happened keeps the common "empty or fully-synced" case
// from adding noise to every command.
func (r *root) reportSync(rep outbox.Report, err error, notice string) {
	out := r.stderr()
	if notice != "" {
		fmt.Fprint(out, notice)
	}
	// Report each outcome independently: a partial run that synced some ops and
	// left others must tell the user about both, not just the leftover.
	if rep.Synced > 0 {
		fmt.Fprintf(out, "synced %d queued %s\n", rep.Synced, pluralEntries(rep.Synced))
	}
	if rep.Failed > 0 {
		fmt.Fprintf(out, "%d queued %s could not be synced — see 'tl sync --list'\n",
			rep.Failed, pluralEntries(rep.Failed))
	}
	if rep.Remaining > 0 {
		fmt.Fprintf(out, "sync incomplete — %s still queued\n", pluralOps(rep.Remaining))
	}
	// ErrLocked simply means another tl is already syncing — not worth reporting.
	if err != nil && !errors.Is(err, outbox.ErrLocked) {
		fmt.Fprintf(out, "sync error: %v\n", err)
	}
}

func pluralEntries(n int) string {
	if n == 1 {
		return "entry"
	}
	return "entries"
}

// pluralOps renders an op count with the right noun ("1 op" / "3 ops").
func pluralOps(n int) string {
	if n == 1 {
		return "1 op"
	}
	return fmt.Sprintf("%d ops", n)
}

// resolveCompany resolves the target company from the --company flag (wins) or
// default_company in config, failing with the standard guidance when neither is set.
func (r *root) resolveCompany(cfg config.Config) (uint64, error) {
	if r.company != 0 {
		return r.company, nil
	}
	if cfg.DefaultCompanyID != 0 {
		return cfg.DefaultCompanyID, nil
	}
	return 0, errors.New(
		"no company set — pass --company <id> or run 'tl config set default_company <id>' (see 'tl companies')")
}

func (r *root) resolveBaseURL(cfg config.Config) string {
	if r.baseURL != "" {
		return r.baseURL
	}
	if cfg.BaseURL != "" {
		return cfg.BaseURL
	}
	return defaultBaseURL
}

func (r *root) client(cfg config.Config) (api.Client, error) {
	opts, err := r.tlsOptions(cfg)
	if err != nil {
		return nil, err
	}
	return api.New(r.resolveBaseURL(cfg), auth.AccessToken, opts...), nil
}

// tlsOptions builds the client's TLS options from flags (which win) and config.
// Returns no options when neither a CA cert nor insecure mode is requested, so
// plain-HTTP base URLs are unaffected.
func (r *root) tlsOptions(cfg config.Config) ([]api.Option, error) {
	insecure := r.insecure || cfg.Insecure

	caPath := r.caCert
	if caPath == "" {
		caPath = cfg.CACertPath
	}

	if !insecure && caPath == "" {
		return nil, nil
	}

	tlsCfg := &tls.Config{InsecureSkipVerify: insecure}

	if caPath != "" {
		pem, err := os.ReadFile(caPath)
		if err != nil {
			return nil, fmt.Errorf("reading ca_cert: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("ca_cert %s: no valid certificates found", caPath)
		}
		tlsCfg.RootCAs = pool
	}

	return []api.Option{api.WithTLSConfig(tlsCfg)}, nil
}
