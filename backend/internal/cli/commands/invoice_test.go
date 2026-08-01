package commands

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// invoiceServer scripts POST /api/invoices/generate. body is served on success;
// disposition (when set) becomes the Content-Disposition header; errMsg with a
// status turns the response into an API error.
type invoiceServer struct {
	body        []byte
	disposition string
	status      int
	errMsg      string

	gotQuery url.Values
}

func (s *invoiceServer) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/invoices/generate", func(w http.ResponseWriter, r *http.Request) {
		s.gotQuery = r.URL.Query()
		if s.errMsg != "" {
			writeJSON(w, s.status, map[string]string{"error": s.errMsg})
			return
		}
		if s.disposition != "" {
			w.Header().Set("Content-Disposition", s.disposition)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(s.body)
	})
	return mux
}

func fakePDF() []byte { return []byte("%PDF-1.7 fake invoice bytes") }

func TestInvoiceWritesServerFilenameIntoDir(t *testing.T) {
	srv := &invoiceServer{body: fakePDF(), disposition: `attachment; filename="invoice-42.pdf"`}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	out, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", dir)
	if err != nil {
		t.Fatalf("invoice returned error: %v", err)
	}

	want := filepath.Join(dir, "invoice-42.pdf")
	got, err := os.ReadFile(want)
	if err != nil {
		t.Fatalf("reading written invoice: %v", err)
	}
	if !bytes.Equal(got, fakePDF()) {
		t.Errorf("written bytes = %q, want %q", got, fakePDF())
	}
	if !strings.Contains(out, want) {
		t.Errorf("stdout = %q, want it to contain %q", out, want)
	}

	if q := srv.gotQuery; q.Get("start") != "2026-06-01" || q.Get("end") != "2026-06-30" ||
		q.Get("company_id") != "1" || q.Get("format") != "pdf" {
		t.Errorf("query = %v, want month 2026-06 expanded for company 1 pdf", srv.gotQuery)
	}
}

func TestInvoiceSanitizesServerFilename(t *testing.T) {
	// A server-supplied filename with directory components must not escape -o.
	srv := &invoiceServer{body: fakePDF(), disposition: `attachment; filename="../../escape.pdf"`}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	if _, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", dir); err != nil {
		t.Fatalf("invoice returned error: %v", err)
	}

	// The document lands inside dir under the basename, not two levels up.
	if _, err := os.ReadFile(filepath.Join(dir, "escape.pdf")); err != nil {
		t.Fatalf("expected sanitized file inside -o dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "..", "..", "escape.pdf")); err == nil {
		t.Fatal("filename traversal escaped the output directory")
	}
}

func TestInvoiceFallbackFilename(t *testing.T) {
	srv := &invoiceServer{body: fakePDF()} // no Content-Disposition
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	if _, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", dir); err != nil {
		t.Fatalf("invoice returned error: %v", err)
	}

	want := filepath.Join(dir, "invoice-2026-06-01-2026-06-30.pdf")
	if _, err := os.ReadFile(want); err != nil {
		t.Fatalf("expected fallback filename %s: %v", want, err)
	}
}

func TestInvoiceExcelFallbackExtension(t *testing.T) {
	srv := &invoiceServer{body: []byte("xlsx-bytes")}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	if _, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "--format", "excel", "-o", dir); err != nil {
		t.Fatalf("invoice returned error: %v", err)
	}
	if _, err := os.ReadFile(filepath.Join(dir, "invoice-2026-06-01-2026-06-30.xlsx")); err != nil {
		t.Fatalf("expected .xlsx fallback filename: %v", err)
	}
	if srv.gotQuery.Get("format") != "excel" {
		t.Errorf("format query = %q, want excel", srv.gotQuery.Get("format"))
	}
}

func TestInvoiceFullOutputPath(t *testing.T) {
	srv := &invoiceServer{body: fakePDF(), disposition: `attachment; filename="ignored.pdf"`}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	target := filepath.Join(t.TempDir(), "custom.pdf")
	if _, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", target); err != nil {
		t.Fatalf("invoice returned error: %v", err)
	}
	if _, err := os.ReadFile(target); err != nil {
		t.Fatalf("expected file at explicit path %s: %v", target, err)
	}
}

func TestInvoiceRefuseOverwriteThenForce(t *testing.T) {
	srv := &invoiceServer{body: fakePDF(), disposition: `attachment; filename="invoice-42.pdf"`}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "invoice-42.pdf")
	if err := os.WriteFile(target, []byte("original"), 0o644); err != nil {
		t.Fatalf("seeding existing file: %v", err)
	}

	_, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", dir)
	if err == nil {
		t.Fatal("expected a refuse-overwrite error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q, want it to mention the existing file", err.Error())
	}
	if got, _ := os.ReadFile(target); string(got) != "original" {
		t.Errorf("file was modified without --force: %q", got)
	}

	if _, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", dir, "--force"); err != nil {
		t.Fatalf("invoice --force returned error: %v", err)
	}
	if got, _ := os.ReadFile(target); !bytes.Equal(got, fakePDF()) {
		t.Errorf("--force did not overwrite: %q", got)
	}
}

func TestInvoiceAPIErrorWritesNoFile(t *testing.T) {
	srv := &invoiceServer{status: http.StatusUnprocessableEntity, errMsg: "no entries in period"}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	_, err := runAuthedCmd(t, ts.URL, "", "invoice", "--month", "2026-06", "-o", dir)
	if err == nil {
		t.Fatal("expected the API error to surface")
	}
	if !strings.Contains(err.Error(), "no entries in period") {
		t.Errorf("error = %q, want the server's message", err.Error())
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("no file should be written on API error, found %d", len(entries))
	}
}

func TestInvoiceMonthAndRangeMutuallyExclusive(t *testing.T) {
	srv := &invoiceServer{body: fakePDF()}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "invoice",
		"--month", "2026-06", "--start", "2026-06-01", "--end", "2026-06-30", "-o", t.TempDir())
	if err == nil {
		t.Fatal("expected a mutual-exclusion error")
	}
	if srv.gotQuery != nil {
		t.Error("no request should be made when the flags conflict")
	}
}

func TestInvoiceStartRequiresEnd(t *testing.T) {
	srv := &invoiceServer{body: fakePDF()}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "invoice", "--start", "2026-06-01", "-o", t.TempDir())
	if err == nil {
		t.Fatal("expected an error when --end is missing")
	}
}

func TestInvoiceExplicitRangePassthrough(t *testing.T) {
	srv := &invoiceServer{body: fakePDF()}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	dir := t.TempDir()
	if _, err := runAuthedCmd(t, ts.URL, "", "invoice",
		"--start", "2026-05-10", "--end", "2026-05-20", "-o", dir); err != nil {
		t.Fatalf("invoice returned error: %v", err)
	}
	if q := srv.gotQuery; q.Get("start") != "2026-05-10" || q.Get("end") != "2026-05-20" {
		t.Errorf("query = %v, want the explicit range passed through", srv.gotQuery)
	}
}

func TestInvoiceRejectsReversedRange(t *testing.T) {
	srv := &invoiceServer{body: fakePDF()}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "invoice",
		"--start", "2026-05-20", "--end", "2026-05-10", "-o", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "before") {
		t.Fatalf("error = %v, want a reversed-range error mentioning 'before'", err)
	}
	if srv.gotQuery != nil {
		t.Error("no request should be made for a reversed range")
	}
}

func TestInvoiceRejectsMalformedRangeDate(t *testing.T) {
	srv := &invoiceServer{body: fakePDF()}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "invoice",
		"--start", "not-a-date", "--end", "2026-05-10", "-o", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "invalid --start") {
		t.Fatalf("error = %v, want an 'invalid --start' error", err)
	}
	if srv.gotQuery != nil {
		t.Error("no request should be made for a malformed date")
	}
}
