package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rescoot/sunshine-cli/internal/config"
)

type TokenSet struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scopes       string    `json:"scopes,omitempty"`
}

func (t *TokenSet) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return false // No expiry set
	}
	return time.Now().After(t.ExpiresAt.Add(-30 * time.Second)) // 30s buffer
}

func tokensPath() string {
	return filepath.Join(config.Dir(), "tokens.json")
}

func LoadTokens() (*TokenSet, error) {
	data, err := os.ReadFile(tokensPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading tokens: %w", err)
	}

	var tokens TokenSet
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, fmt.Errorf("parsing tokens: %w", err)
	}

	return &tokens, nil
}

func SaveTokens(tokens *TokenSet) error {
	if err := os.MkdirAll(config.Dir(), 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling tokens: %w", err)
	}

	return os.WriteFile(tokensPath(), data, 0o600)
}

func ClearTokens() error {
	path := tokensPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}
