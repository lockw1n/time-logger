package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const requestTimeout = 10 * time.Second

// httpClient is the Client implementation that talks HTTP+JSON to the backend
// through nginx.
type httpClient struct {
	http    *http.Client
	baseURL string
	tokenFn func() (string, error)
}

// Option customizes the client at construction time.
type Option func(*httpClient)

// WithTLSConfig makes the client use the given TLS configuration — needed to
// trust a self-signed backend cert (custom RootCAs) or skip verification.
// A nil config is ignored.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(c *httpClient) {
		if cfg == nil {
			return
		}
		tr := http.DefaultTransport.(*http.Transport).Clone()
		tr.TLSClientConfig = cfg
		c.http.Transport = tr
	}
}

func New(baseURL string, tokenFn func() (string, error), opts ...Option) Client {
	return newHTTPClient(baseURL, tokenFn, opts...)
}

func newHTTPClient(baseURL string, tokenFn func() (string, error), opts ...Option) *httpClient {
	c := &httpClient{
		http:    &http.Client{Timeout: requestTimeout},
		baseURL: strings.TrimRight(baseURL, "/"),
		tokenFn: tokenFn,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *httpClient) Login(ctx context.Context, email, password string) (LoginResult, error) {
	var result LoginResult

	err := c.do(ctx, http.MethodPost, "/api/auth/login", loginRequest{Email: email, Password: password}, &result, false)

	return result, err
}

func (c *httpClient) Health(ctx context.Context) error {
	return c.do(ctx, http.MethodGet, "/health", nil, nil, false)
}

func (c *httpClient) Companies(ctx context.Context) ([]Company, error) {
	var out []Company
	err := c.do(ctx, http.MethodGet, CompaniesPath(), nil, &out, true)
	return out, err
}

func (c *httpClient) Activities(ctx context.Context, companyID uint64) ([]Activity, error) {
	var out []Activity
	err := c.do(ctx, http.MethodGet, ActivitiesPath(companyID), nil, &out, true)
	return out, err
}

func (c *httpClient) Timesheet(ctx context.Context, companyID uint64, start, end string) (TimesheetResponse, error) {
	var out TimesheetResponse
	err := c.do(ctx, http.MethodGet, TimesheetPath(companyID, start, end), nil, &out, true)
	normalizeTimesheetDates(&out)
	return out, err
}

// GetRaw performs an authenticated GET and returns the response body verbatim,
// unparsed and un-normalized — this is what --json prints, so scripts see the
// backend's true payload rather than the CLI's typed/normalized view.
func (c *httpClient) GetRaw(ctx context.Context, path string) (json.RawMessage, error) {
	var out json.RawMessage
	err := c.do(ctx, http.MethodGet, path, nil, &out, true)
	return out, err
}

// Path builders keep query construction in one place, shared by the typed
// methods above and the raw GET used for --json.

func CompaniesPath() string { return "/api/companies" }

func ActivitiesPath(companyID uint64) string {
	q := url.Values{"company_id": {strconv.FormatUint(companyID, 10)}}
	return "/api/activities?" + q.Encode()
}

func TimesheetPath(companyID uint64, start, end string) string {
	q := url.Values{
		"company_id": {strconv.FormatUint(companyID, 10)},
		"start":      {start},
		"end":        {end},
	}
	return "/api/timesheet?" + q.Encode()
}

func (c *httpClient) CreateEntry(ctx context.Context, req CreateEntryRequest) (Entry, error) {
	var out Entry
	err := c.do(ctx, http.MethodPost, "/api/entries", req, &out, true)
	out.Date = canonicalDate(out.Date)
	return out, err
}

func (c *httpClient) UpdateEntry(ctx context.Context, id uint64, req UpdateEntryRequest) (Entry, error) {
	var out Entry
	err := c.do(ctx, http.MethodPut, entryPath(id), req, &out, true)
	out.Date = canonicalDate(out.Date)
	return out, err
}

func (c *httpClient) GetEntry(ctx context.Context, id uint64) (Entry, error) {
	var out Entry
	err := c.do(ctx, http.MethodGet, entryPath(id), nil, &out, true)
	out.Date = canonicalDate(out.Date)
	return out, err
}

func (c *httpClient) DeleteEntry(ctx context.Context, id uint64) error {
	return c.do(ctx, http.MethodDelete, entryPath(id), nil, nil, true)
}

func entryPath(id uint64) string {
	return "/api/entries/" + strconv.FormatUint(id, 10)
}

// invoiceTimeout overrides the client's default 10s deadline for invoice
// generation only: the request round-trips a headless PDF renderer server-side
// and is legitimately slow.
const invoiceTimeout = 60 * time.Second

// GenerateInvoice POSTs to /api/invoices/generate (params are query, body is
// empty) and returns the document bytes with a resolved output filename. It uses
// its own longer-deadline client so the slow render doesn't trip the 10s default.
func (c *httpClient) GenerateInvoice(ctx context.Context, req GenerateInvoiceRequest) (InvoiceFile, error) {
	endpoint := c.baseURL + InvoicePath(req.CompanyID, req.Start, req.End, req.Format)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return InvoiceFile{}, fmt.Errorf("building request: %w", err)
	}

	token, err := c.tokenFn()
	if err != nil {
		return InvoiceFile{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: invoiceTimeout, Transport: c.http.Transport}
	resp, err := client.Do(httpReq)
	if err != nil {
		return InvoiceFile{}, transportError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return InvoiceFile{}, responseError(resp, true)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return InvoiceFile{}, fmt.Errorf("reading invoice: %w", err)
	}

	return InvoiceFile{
		Filename: invoiceFilename(resp.Header.Get("Content-Disposition"), req),
		Data:     data,
	}, nil
}

// InvoicePath builds the invoice-generation URL. All parameters are query
// params despite the POST method (the backend binds them from the query string).
func InvoicePath(companyID uint64, start, end, format string) string {
	q := url.Values{
		"company_id": {strconv.FormatUint(companyID, 10)},
		"start":      {start},
		"end":        {end},
		"format":     {format},
	}
	return "/api/invoices/generate?" + q.Encode()
}

// invoiceFilename prefers the server's Content-Disposition filename and falls
// back to invoice-<start>-<end>.<ext> when the header is absent or unparseable.
func invoiceFilename(disposition string, req GenerateInvoiceRequest) string {
	if disposition != "" {
		if _, params, err := mime.ParseMediaType(disposition); err == nil {
			// filepath.Base strips any directory components a server might smuggle
			// in via the filename, so the document can never be written outside the
			// caller's chosen output location.
			if name := params["filename"]; name != "" {
				if base := filepath.Base(name); base != "." && base != ".." && base != string(filepath.Separator) {
					return base
				}
			}
		}
	}
	return fmt.Sprintf("invoice-%s-%s.%s", req.Start, req.End, invoiceExt(req.Format))
}

func invoiceExt(format string) string {
	if format == "excel" {
		return "xlsx"
	}
	return "pdf"
}

func (c *httpClient) do(ctx context.Context, method, path string, body, out any, authed bool) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encoding request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if authed {
		token, err := c.tokenFn()
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return transportError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return responseError(resp, authed)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// transportError classifies a failed *http.Client.Do into one of two sentinels.
// A cancelled or timed-out request (context.Canceled/DeadlineExceeded — Ctrl-C,
// or the client's own timeout after the request was sent) is ambiguous: the
// server may already have committed it, so it is ErrInterrupted and must not be
// auto-queued. Everything else (connection refused, DNS, no route) means the
// request never landed, so it is ErrUnreachable and safe to queue and replay.
func transportError(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", ErrInterrupted, err)
	}
	return fmt.Errorf("%w: %v", ErrUnreachable, err)
}

func responseError(resp *http.Response, authed bool) error {
	if authed && resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	msg := strings.TrimSpace(string(data))

	var payload errorResponse
	if err := json.Unmarshal(data, &payload); err == nil && payload.Error != "" {
		msg = payload.Error
	}

	if msg == "" {
		msg = http.StatusText(resp.StatusCode)
	}

	return &APIError{Status: resp.StatusCode, Message: msg}
}
