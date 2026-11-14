package domain

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	CodeDigits      = 6
	CodeLifetime    = 10 * time.Minute
	CodeMaxAttempts = 5
)

var codeUpperBound = new(big.Int).Exp(big.NewInt(10), big.NewInt(CodeDigits), nil)

// LoginCode is the stored side of a code: its hash, never the code.
type LoginCode struct {
	UserID    int64
	CodeHash  string
	ExpiresAt time.Time
	Attempts  int
}

type CodeIssue struct {
	UserID int64
	Now    time.Time
}

// IssuedCode carries the code in clear only so it can be sent.
type IssuedCode struct {
	Code      string
	LoginCode LoginCode
}

func IssueLoginCode(random io.Reader, issue CodeIssue) (IssuedCode, error) {
	code, err := randomDigits(random)
	if err != nil {
		return IssuedCode{}, err
	}
	stored := LoginCode{UserID: issue.UserID, CodeHash: HashSecret(code), ExpiresAt: issue.Now.Add(CodeLifetime)}
	return IssuedCode{Code: code, LoginCode: stored}, nil
}

func randomDigits(random io.Reader) (string, error) {
	n, err := rand.Int(random, codeUpperBound)
	if err != nil {
		return "", fmt.Errorf("generate login code: %w", err)
	}
	return fmt.Sprintf("%0*d", CodeDigits, n), nil
}

// ParseCode checks the code at the edge: exactly six ASCII digits after
// trimming. A code of the wrong shape never reaches verification, so it does
// not use up an attempt.
func ParseCode(raw string) (string, error) {
	code := strings.TrimSpace(raw)
	if len(code) != CodeDigits || strings.ContainsFunc(code, isNotASCIIDigit) {
		return "", shared.Invalid(KeyCodeFormatInvalid, nil)
	}
	return code, nil
}

func isNotASCIIDigit(char rune) bool { return char < '0' || char > '9' }

func (c LoginCode) Matches(entered string) bool { return secretMatches(entered, c.CodeHash) }

func (c LoginCode) IsExpiredAt(now time.Time) bool { return hasExpired(c.ExpiresAt, now) }

// AttemptsExhausted is true once the failed attempts reach the limit: that code no longer opens.
func AttemptsExhausted(attempts int) bool { return attempts >= CodeMaxAttempts }
