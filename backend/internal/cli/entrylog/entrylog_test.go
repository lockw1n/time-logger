package entrylog

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// fakeClient is a minimal scriptable api.Client for the submit/merge tests.
type fakeClient struct {
	createErr error
	ts        api.TimesheetResponse
	tsErr     error
	updateErr error
}

func (f *fakeClient) Login(context.Context, string, string) (api.LoginResult, error) {
	return api.LoginResult{}, nil
}
func (f *fakeClient) Health(context.Context) error { return nil }
func (f *fakeClient) Companies(context.Context) ([]api.Company, error) {
	return nil, nil
}
func (f *fakeClient) Activities(context.Context, uint64) ([]api.Activity, error) {
	return nil, nil
}
func (f *fakeClient) Timesheet(context.Context, uint64, string, string) (api.TimesheetResponse, error) {
	return f.ts, f.tsErr
}
func (f *fakeClient) CreateEntry(context.Context, api.CreateEntryRequest) (api.Entry, error) {
	return api.Entry{}, f.createErr
}
func (f *fakeClient) UpdateEntry(_ context.Context, id uint64, _ api.UpdateEntryRequest) (api.Entry, error) {
	if f.updateErr != nil {
		return api.Entry{}, f.updateErr
	}
	return api.Entry{ID: id}, nil
}
func (f *fakeClient) GetEntry(_ context.Context, id uint64) (api.Entry, error) {
	return api.Entry{ID: id}, nil
}
func (f *fakeClient) DeleteEntry(context.Context, uint64) error { return nil }
func (f *fakeClient) GenerateInvoice(context.Context, api.GenerateInvoiceRequest) (api.InvoiceFile, error) {
	return api.InvoiceFile{}, nil
}
func (f *fakeClient) GetRaw(context.Context, string) (json.RawMessage, error) { return nil, nil }

func testRequest() Request {
	return Request{
		CompanyID:  1,
		TicketCode: "APP-1",
		Activity:   api.Activity{ID: 5, Name: "dev"},
		Date:       "2026-08-01",
		Minutes:    30,
	}
}

// A create that never reaches the backend is safe to queue: Submit must tag it
// with ErrCreateUnreachable (and keep api.ErrUnreachable in the chain so the
// replay classifier still treats it as retryable).
func TestSubmitCreateUnreachableIsQueueable(t *testing.T) {
	client := &fakeClient{createErr: api.ErrUnreachable}

	err := Submit(context.Background(), client, testRequest(), Options{Out: io.Discard, SkipConfirm: true})

	if !errors.Is(err, ErrCreateUnreachable) {
		t.Errorf("create-unreachable error = %v, want it to wrap ErrCreateUnreachable", err)
	}
	if !errors.Is(err, api.ErrUnreachable) {
		t.Errorf("create-unreachable error = %v, want it to still wrap api.ErrUnreachable", err)
	}
}

// Regression for the double-count finding: when the create hits a 409 and the
// merge's UpdateEntry then fails unreachable, the update may already have
// committed. Submit must NOT tag it ErrCreateUnreachable (so the caller won't
// re-queue it as a fresh create), while keeping it retryable for replay.
func TestSubmitMergeUpdateUnreachableIsNotQueueable(t *testing.T) {
	client := &fakeClient{
		createErr: &api.APIError{Status: http.StatusConflict, Message: "duplicate"},
		ts: api.TimesheetResponse{
			Rows: []api.TimesheetRow{{
				Ticket:   api.TimesheetTicket{Code: "APP-1"},
				Activity: api.TimesheetActivity{ID: 5},
				Entries:  []api.TimesheetEntry{{ID: 42, Date: "2026-08-01", DurationMinutes: 60}},
			}},
		},
		updateErr: api.ErrUnreachable,
	}

	err := Submit(context.Background(), client, testRequest(), Options{Out: io.Discard, SkipConfirm: true})

	if err == nil {
		t.Fatal("expected an error from a failed merge update")
	}
	if errors.Is(err, ErrCreateUnreachable) {
		t.Errorf("merge-update error = %v, must NOT be tagged ErrCreateUnreachable (would double-count on re-queue)", err)
	}
	if !errors.Is(err, api.ErrUnreachable) {
		t.Errorf("merge-update error = %v, want it to remain retryable (wrap api.ErrUnreachable)", err)
	}
}
