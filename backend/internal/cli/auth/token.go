package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/config"
)

const (
	tokenFile = "token.json"

	// expiryMargin treats tokens about to expire as already expired so a
	// command doesn't start with a token that dies mid-request.
	expiryMargin = 30 * time.Second
)

var (
	ErrNotLoggedIn  = errors.New("not logged in")
	ErrNoToken      = fmt.Errorf("%w: no token", ErrNotLoggedIn)
	ErrTokenExpired = fmt.Errorf("%w: token expired", ErrNotLoggedIn)
)

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func Load() (Token, error) {
	path, err := tokenPath()
	if err != nil {
		return Token{}, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Token{}, ErrNoToken
	}
	if err != nil {
		return Token{}, fmt.Errorf("reading %s: %w", tokenFile, err)
	}

	var tok Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return Token{}, fmt.Errorf("parsing %s: %w", tokenFile, err)
	}

	if tok.AccessToken == "" {
		return Token{}, ErrNoToken
	}

	if !time.Now().Add(expiryMargin).Before(tok.ExpiresAt) {
		return Token{}, ErrTokenExpired
	}

	return tok, nil
}

func Save(tok Token) error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding %s: %w", tokenFile, err)
	}

	return config.WriteFileAtomic(path, data, 0o600)
}

func Delete() error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("removing %s: %w", tokenFile, err)
	}

	return nil
}

// AccessToken is the token provider passed to api.New.
func AccessToken() (string, error) {
	tok, err := Load()
	if err != nil {
		return "", err
	}

	return tok.AccessToken, nil
}

func tokenPath() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, tokenFile), nil
}
