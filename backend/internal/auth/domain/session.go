package domain

import (
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

const (
	SessionLifetime   = 30 * 24 * time.Hour
	SessionTokenBytes = 32
)

// Session is the stored side of a session: the hash of its token, never the token.
type Session struct {
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
}

type SessionOpening struct {
	UserID int64
	Now    time.Time
}

// OpenedSession carries the token in clear only to hand it to the browser once.
type OpenedSession struct {
	Token   string
	Session Session
}

func NewSession(random io.Reader, opening SessionOpening) (OpenedSession, error) {
	bytes, err := readRandom(random, SessionTokenBytes)
	if err != nil {
		return OpenedSession{}, fmt.Errorf("generate session token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	stored := Session{UserID: opening.UserID, TokenHash: HashSecret(token), ExpiresAt: opening.Now.Add(SessionLifetime)}
	return OpenedSession{Token: token, Session: stored}, nil
}

func (s Session) IsExpiredAt(now time.Time) bool { return hasExpired(s.ExpiresAt, now) }
