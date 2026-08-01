package tui

import (
	"context"
	"encoding/json"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// fakeClient is a scriptable api.Client for the update tests. It records the
// writes it receives and returns pre-set responses, so a Cmd's effect can be
// asserted without a real backend.
type fakeClient struct {
	ts         api.TimesheetResponse
	companies  []api.Company
	activities []api.Activity

	tsErr     error
	updateErr error
	deleteErr error
	createErr error

	updated []updateCall
	deleted []uint64
	created []api.CreateEntryRequest
}

type updateCall struct {
	id  uint64
	req api.UpdateEntryRequest
}

func (f *fakeClient) Login(context.Context, string, string) (api.LoginResult, error) {
	return api.LoginResult{}, nil
}
func (f *fakeClient) Health(context.Context) error { return nil }

func (f *fakeClient) Companies(context.Context) ([]api.Company, error) {
	return f.companies, nil
}
func (f *fakeClient) Activities(context.Context, uint64) ([]api.Activity, error) {
	return f.activities, nil
}
func (f *fakeClient) Timesheet(context.Context, uint64, string, string) (api.TimesheetResponse, error) {
	return f.ts, f.tsErr
}
func (f *fakeClient) CreateEntry(_ context.Context, req api.CreateEntryRequest) (api.Entry, error) {
	f.created = append(f.created, req)
	if f.createErr != nil {
		return api.Entry{}, f.createErr
	}
	return api.Entry{ID: 999, DurationMinutes: req.DurationMinutes}, nil
}
func (f *fakeClient) UpdateEntry(_ context.Context, id uint64, req api.UpdateEntryRequest) (api.Entry, error) {
	f.updated = append(f.updated, updateCall{id: id, req: req})
	if f.updateErr != nil {
		return api.Entry{}, f.updateErr
	}
	return api.Entry{ID: id}, nil
}
func (f *fakeClient) GetEntry(_ context.Context, id uint64) (api.Entry, error) {
	return api.Entry{ID: id}, nil
}
func (f *fakeClient) DeleteEntry(_ context.Context, id uint64) error {
	f.deleted = append(f.deleted, id)
	return f.deleteErr
}
func (f *fakeClient) GenerateInvoice(context.Context, api.GenerateInvoiceRequest) (api.InvoiceFile, error) {
	return api.InvoiceFile{}, nil
}
func (f *fakeClient) GetRaw(context.Context, string) (json.RawMessage, error) {
	return nil, nil
}
