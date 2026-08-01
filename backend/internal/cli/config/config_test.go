package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirEnvOverride(t *testing.T) {
	t.Setenv(dirEnv, "/custom/tl")

	dir, err := Dir()
	if err != nil {
		t.Fatalf("Dir() error: %v", err)
	}
	if dir != "/custom/tl" {
		t.Fatalf("Dir() = %q, want %q", dir, "/custom/tl")
	}
}

func TestDirDefault(t *testing.T) {
	t.Setenv(dirEnv, "")
	t.Setenv("HOME", "/home/test")

	dir, err := Dir()
	if err != nil {
		t.Fatalf("Dir() error: %v", err)
	}
	if want := filepath.Join("/home/test", ".config", "tl"); dir != want {
		t.Fatalf("Dir() = %q, want %q", dir, want)
	}
}

func TestLoadMissingFile(t *testing.T) {
	t.Setenv(dirEnv, t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg != (Config{}) {
		t.Fatalf("Load() = %+v, want zero value", cfg)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Setenv(dirEnv, t.TempDir())

	want := Config{
		BaseURL:          "http://localhost:9999",
		DefaultCompanyID: 42,
		Email:            "user@example.com",
	}

	if err := Save(want); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "tl")
	t.Setenv(dirEnv, dir)

	if err := Save(Config{Email: "user@example.com"}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat config dir: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Fatalf("config dir mode = %o, want 0700", perm)
	}
}
