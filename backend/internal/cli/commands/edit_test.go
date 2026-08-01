package commands

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEditRequiresMutationFlag(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "edit", "421")
	if err == nil {
		t.Fatal("expected an error when no mutation flag is given")
	}
	var coder interface{ ExitCode() int }
	if !errors.As(err, &coder) || coder.ExitCode() != 1 {
		t.Fatalf("expected exit code 1, got %v", err)
	}
	if srv.putCalled {
		t.Error("no PUT should be sent when there is nothing to change")
	}
}

func TestEditDurationPutsAndConfirms(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	out, err := runAuthedCmd(t, ts.URL, "", "edit", "421", "--duration", "2h", "--yes")
	if err != nil {
		t.Fatalf("edit returned error: %v", err)
	}
	if !srv.putCalled {
		t.Fatal("expected a PUT, none was sent")
	}

	var put struct {
		DurationMinutes *int    `json:"duration_minutes"`
		Comment         *string `json:"comment"`
	}
	if err := json.Unmarshal(srv.putBody, &put); err != nil {
		t.Fatalf("decoding PUT body: %v", err)
	}
	if put.DurationMinutes == nil || *put.DurationMinutes != 120 {
		t.Errorf("PUT duration = %v, want 120", put.DurationMinutes)
	}
	if put.Comment != nil {
		t.Errorf("PUT comment = %v, want nil (unchanged)", put.Comment)
	}
	if !strings.Contains(out, "#421") {
		t.Errorf("output %q should reference entry #421", out)
	}
}

func TestEditAbortsOnNonInteractiveStdin(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "edit", "421", "--duration", "2h") // no --yes, EOF stdin
	if err == nil {
		t.Fatal("expected the edit to abort on EOF, got nil")
	}
	if srv.putCalled {
		t.Error("no PUT should be sent when the confirmation is not accepted")
	}
}

func TestEditConfirmedViaStdin(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "y\n", "edit", "421", "--duration", "2h")
	if err != nil {
		t.Fatalf("edit returned error: %v", err)
	}
	if !srv.putCalled {
		t.Error("expected a PUT after the prompt was accepted")
	}
}

func TestEditClearCommentRejected(t *testing.T) {
	srv := &entryServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "edit", "421", "--comment", "  ", "--yes")
	if err == nil {
		t.Fatal("expected an error for an empty comment")
	}
	if !strings.Contains(err.Error(), "can't be cleared") {
		t.Errorf("error = %q, want it to explain comments can't be cleared", err.Error())
	}
	if srv.putCalled {
		t.Error("no PUT should be sent for a rejected clear-comment")
	}
}

func TestEditNotFound(t *testing.T) {
	srv := &entryServer{notFound: true}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	_, err := runAuthedCmd(t, ts.URL, "", "edit", "421", "--duration", "2h", "--yes")
	if err == nil {
		t.Fatal("expected a not-found error")
	}
	if !strings.Contains(err.Error(), "entry 421 not found") {
		t.Errorf("error = %q, want 'entry 421 not found'", err.Error())
	}
}
