package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func noToken() (string, error) {
	return "", errors.New("no token expected in this test")
}

func TestLoginSuccess(t *testing.T) {
	expiresAt := time.Now().Add(30 * time.Minute).UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/auth/login" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decoding login request: %v", err)
		}
		if req.Email != "user@example.com" || req.Password != "secret" {
			t.Errorf("login request = %+v", req)
		}

		json.NewEncoder(w).Encode(LoginResult{AccessToken: "jwt-token", ExpiresAt: expiresAt})
	}))
	defer server.Close()

	result, err := New(server.URL, noToken).Login(context.Background(), "user@example.com", "secret")
	if err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if result.AccessToken != "jwt-token" {
		t.Fatalf("AccessToken = %q, want %q", result.AccessToken, "jwt-token")
	}
	if !result.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("ExpiresAt = %v, want %v", result.ExpiresAt, expiresAt)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
	}))
	defer server.Close()

	_, err := New(server.URL, noToken).Login(context.Background(), "user@example.com", "wrong")

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Login() error = %v, want *APIError", err)
	}
	if apiErr.Status != http.StatusUnauthorized || apiErr.Message != "invalid credentials" {
		t.Fatalf("APIError = %+v", apiErr)
	}
	// Login is unauthenticated; its 401 must not map to the session sentinel.
	if errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Login() 401 must not be ErrUnauthorized")
	}
}

func TestNonJSONErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("bad gateway"))
	}))
	defer server.Close()

	err := New(server.URL, noToken).Health(context.Background())

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Health() error = %v, want *APIError", err)
	}
	if apiErr.Status != http.StatusBadGateway || apiErr.Message != "bad gateway" {
		t.Fatalf("APIError = %+v", apiErr)
	}
}

func TestUnreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	server.Close()

	err := New(server.URL, noToken).Health(context.Background())
	if !errors.Is(err, ErrUnreachable) {
		t.Fatalf("Health() error = %v, want ErrUnreachable", err)
	}
	// A refused connection never reached the server, so it must NOT masquerade as
	// an interrupted (maybe-delivered) request.
	if errors.Is(err, ErrInterrupted) {
		t.Fatalf("connection refused wrongly classified as ErrInterrupted: %v", err)
	}
}

// A cancelled request may have reached the server, so it must surface as
// ErrInterrupted (not ErrUnreachable) — the outbox must never auto-queue it, or
// a create the server already committed would be replayed and double-counted.
func TestCancelledRequestIsInterruptedNotUnreachable(t *testing.T) {
	started := make(chan struct{})
	serverDone := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-r.Context().Done():
		case <-serverDone:
		}
	}))
	defer server.Close()
	defer close(serverDone)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- New(server.URL, noToken).Health(ctx)
	}()

	<-started
	cancel()

	err := <-errCh
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("cancelled request error = %v, want ErrInterrupted", err)
	}
	if errors.Is(err, ErrUnreachable) {
		t.Fatalf("cancelled request wrongly classified as ErrUnreachable: %v", err)
	}
}

func TestBearerHeaderAttached(t *testing.T) {
	var gotAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	client := newHTTPClient(server.URL, func() (string, error) { return "jwt-token", nil })
	if err := client.do(context.Background(), http.MethodGet, "/api/companies", nil, nil, true); err != nil {
		t.Fatalf("do() error: %v", err)
	}

	if gotAuth != "Bearer jwt-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer jwt-token")
	}
}

func TestAuthenticatedUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid or expired token"})
	}))
	defer server.Close()

	client := newHTTPClient(server.URL, func() (string, error) { return "stale", nil })
	err := client.do(context.Background(), http.MethodGet, "/api/companies", nil, nil, true)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("do() error = %v, want ErrUnauthorized", err)
	}
}

func TestTokenProviderErrorPropagates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("request must not be sent when the token provider fails")
	}))
	defer server.Close()

	wantErr := errors.New("no token")
	client := newHTTPClient(server.URL, func() (string, error) { return "", wantErr })

	err := client.do(context.Background(), http.MethodGet, "/api/companies", nil, nil, true)
	if !errors.Is(err, wantErr) {
		t.Fatalf("do() error = %v, want %v", err, wantErr)
	}
}
