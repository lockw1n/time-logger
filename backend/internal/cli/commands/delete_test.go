package commands

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteConfirmsAndDeletes(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	out, err := runAuthedCmd(t, ts.URL, "", "delete", "421", "--yes")
	if err != nil {
		t.Fatalf("delete returned error: %v", err)
	}
	if !srv.deleteCalled {
		t.Fatal("expected a DELETE, none was sent")
	}
	if !strings.Contains(out, "deleted entry #421") {
		t.Errorf("output = %q, want it to confirm deletion of #421", out)
	}
}

func TestDeleteAbortsOnNonInteractiveStdin(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "delete", "421") // no --yes, EOF stdin
	if err == nil {
		t.Fatal("expected the delete to abort on EOF, got nil")
	}
	if srv.deleteCalled {
		t.Error("no DELETE should be sent when the confirmation is not accepted")
	}
}

func TestDeleteConfirmedViaStdin(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "y\n", "delete", "421")
	if err != nil {
		t.Fatalf("delete returned error: %v", err)
	}
	if !srv.deleteCalled {
		t.Error("expected a DELETE after the prompt was accepted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv := &entryServer{notFound: true}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "delete", "421", "--yes")
	if err == nil {
		t.Fatal("expected a not-found error")
	}
	if !strings.Contains(err.Error(), "entry 421 not found") {
		t.Errorf("error = %q, want 'entry 421 not found'", err.Error())
	}
	if srv.deleteCalled {
		t.Error("no DELETE should be sent when the entry is missing")
	}
}
