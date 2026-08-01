package auth

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		token   *Token
		wantErr error
	}{
		{
			name:    "missing file",
			token:   nil,
			wantErr: ErrNoToken,
		},
		{
			name:    "expired token",
			token:   &Token{AccessToken: "tok", ExpiresAt: time.Now().Add(-time.Hour)},
			wantErr: ErrTokenExpired,
		},
		{
			name:    "expires within safety margin",
			token:   &Token{AccessToken: "tok", ExpiresAt: time.Now().Add(10 * time.Second)},
			wantErr: ErrTokenExpired,
		},
		{
			name:  "valid token",
			token: &Token{AccessToken: "tok", ExpiresAt: time.Now().Add(time.Hour)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TL_CONFIG_DIR", t.TempDir())

			if tt.token != nil {
				if err := Save(*tt.token); err != nil {
					t.Fatalf("Save() error: %v", err)
				}
			}

			tok, err := Load()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Load() error = %v, want %v", err, tt.wantErr)
				}
				if !errors.Is(err, ErrNotLoggedIn) {
					t.Fatalf("Load() error = %v, want it to wrap ErrNotLoggedIn", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}
			if tok.AccessToken != tt.token.AccessToken {
				t.Fatalf("Load() token = %q, want %q", tok.AccessToken, tt.token.AccessToken)
			}
		})
	}
}

func TestSaveFileMode(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TL_CONFIG_DIR", dir)

	if err := Save(Token{AccessToken: "tok", ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, tokenFile))
	if err != nil {
		t.Fatalf("stat token file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("token file mode = %o, want 0600", perm)
	}
}

func TestDelete(t *testing.T) {
	t.Setenv("TL_CONFIG_DIR", t.TempDir())

	if err := Save(Token{AccessToken: "tok", ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if err := Delete(); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	if _, err := Load(); !errors.Is(err, ErrNoToken) {
		t.Fatalf("Load() after Delete() error = %v, want ErrNoToken", err)
	}

	if err := Delete(); err != nil {
		t.Fatalf("Delete() on missing file error: %v", err)
	}
}
