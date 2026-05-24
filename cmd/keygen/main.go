package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// secret returns the signing secret from the MK_SECRET environment variable.
// If MK_SECRET is not set, it falls back to a default (for development only).
func secret() string {
	if s := os.Getenv("MK_SECRET"); s != "" {
		return s
	}
	return "mK9vN2xP5rF8jL3wH6tY1qA4cG7uE0bD5sM8pR2vJ9nK6wT3yF1aC4hU7oX0zL5"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: keygen user@email.com")
		return
	}
	email := strings.ToLower(strings.TrimSpace(os.Args[1]))
	mac := hmac.New(sha256.New, []byte(secret()))
	mac.Write([]byte(email))
	sig := hex.EncodeToString(mac.Sum(nil))[:16]
	fmt.Printf("%s-%s\n", email, sig)
	fmt.Println("\nCreate a file named license.key with the above line as its content.")
}
