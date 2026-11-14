package domain

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// HashSecret is how codes and session tokens are stored: never in clear.
func HashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// secretMatches compares in constant time so the time taken does not leak the hash.
func secretMatches(secret, hash string) bool {
	return subtle.ConstantTimeCompare([]byte(HashSecret(secret)), []byte(hash)) == 1
}

func readRandom(random io.Reader, size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := io.ReadFull(random, bytes); err != nil {
		return nil, fmt.Errorf("read random bytes: %w", err)
	}
	return bytes, nil
}

// hasExpired treats the expiry instant itself as already expired.
func hasExpired(expiresAt, now time.Time) bool {
	return !now.Before(expiresAt)
}
