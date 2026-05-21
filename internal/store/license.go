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

const secret = "mK9vN2xP5rF8jL3wH6tY1qA4cG7uE0bD5sM8pR2vJ9nK6wT3yF1aC4hU7oX0zL5"
const trialFile = "trial.start"


func ValidLicense() bool {
    exePath, _ := os.Executable() // added line 19,20,21
    exeDir := filepath.Dir(exePath)
    data, err := os.ReadFile(filepath.Join(exeDir, "license.key"))
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
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
    expected := hex.EncodeToString(mac.Sum(nil))[:16]
    return signature == expected
}

// TrialDaysLeft returns remaining trial days. Negative means expired.
// Returns 30 on first run, decreasing each day

func TrialDaysLeft() int {
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    trialPath := filepath.Join(exeDir, trialFile)
    data, err := os.ReadFile(trialPath)    //data, err := os.ReadFile(trialFile) -- changed from this to the uncommented one
    if err != nil {
        // First run - write today's date
        now := time.Now().Format(time.RFC3339)
        os.WriteFile(trialFile, []byte(now), 0644)
        return 30

    }
    start, err := time.Parse(time.RFC3339, strings.TrimSpace(string(data)))
    if err != nil {
         // Corrupted - reset
         now := time.Now().Format(time.RFC3339)
         os.WriteFile(trialFile, []byte(now), 0644)
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