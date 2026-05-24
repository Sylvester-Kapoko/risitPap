package store

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// secret returns the signing secret from the MK_SECRET environment variable.
// If MK_SECRET is not set, it falls back to a default (for development only).
func secret() string {
	if s := os.Getenv("MK_SECRET"); s != "" {
		return s
	}
	return "mK9vN2xP5rF8jL3wH6tY1qA4cG7uE0bD5sM8pR2vJ9nK6wT3yF1aC4hU7oX0zL5"
}

const trialFile = "trial.start"

func ValidLicense() bool {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	data, err := os.ReadFile(filepath.Join(exeDir, "license.key")) // #nosec G304
	if err != nil {
		return false
	}
	key := strings.TrimSpace(string(data))
	return validate(key)
}

func validate(key string) bool {
	parts := strings.SplitN(key, "-", 2)
	if len(parts) != 2 {
		return false
	}
	email, signature := parts[0], parts[1]
	mac := hmac.New(sha256.New, []byte(secret()))
	mac.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	expected := hex.EncodeToString(mac.Sum(nil))[:16]
	return signature == expected
}

func TrialDaysLeft() int {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	trialPath := filepath.Join(exeDir, trialFile)
	data, err := os.ReadFile(trialPath) // #nosec G304
	if err != nil {
		now := time.Now().Format(time.RFC3339)
		_ = os.WriteFile(trialPath, []byte(now), 0600) // tightened permissions
		return 30
	}
	start, err := time.Parse(time.RFC3339, strings.TrimSpace(string(data)))
	if err != nil {
		now := time.Now().Format(time.RFC3339)
		_ = os.WriteFile(trialPath, []byte(now), 0600) // tightened permissions
		return 30
	}

	elapsed := time.Since(start)
	remaining := 30 - int(elapsed.Hours()/24)
	return remaining
}

// ValidLicense satisfies the ReceiptStore interface.
func (s *Store) ValidLicense() bool {
	return ValidLicense()
}

// TrialDaysLeft satisfies the ReceiptStore interface.
func (s *Store) TrialDaysLeft() int {
	return TrialDaysLeft()
}
