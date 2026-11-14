package domain

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
	"time"
)

func randomOf(fill byte) *bytes.Reader {
	return bytes.NewReader(bytes.Repeat([]byte{fill}, SessionTokenBytes))
}

func TestNewSession_TokenIsThirtyTwoRandomBytesInURLSafeBase64(t *testing.T) {
	got, err := NewSession(randomOf(0xFB), SessionOpening{UserID: 7, Now: testNow})

	if err != nil {
		t.Fatalf("NewSession error = %v, want nil", err)
	}
	decoded, decodeErr := base64.RawURLEncoding.DecodeString(got.Token)
	if decodeErr != nil || !bytes.Equal(decoded, bytes.Repeat([]byte{0xFB}, 32)) {
		t.Errorf("token %q decodes to %x, %v; want the 32 random bytes", got.Token, decoded, decodeErr)
	}
}

func TestNewSession_ExpiresInThirtyDays(t *testing.T) {
	got, err := NewSession(randomOf(1), SessionOpening{UserID: 7, Now: testNow})

	if err != nil || !got.Session.ExpiresAt.Equal(testNow.Add(30*24*time.Hour)) || got.Session.UserID != 7 {
		t.Errorf("session = %+v, %v; want user 7 expiring at %v", got.Session, err, testNow.Add(30*24*time.Hour))
	}
}

func TestNewSession_StoresOnlyTheHashOfTheToken(t *testing.T) {
	got, err := NewSession(randomOf(1), SessionOpening{UserID: 7, Now: testNow})

	if err != nil || got.Session.TokenHash != HashSecret(got.Token) || got.Session.TokenHash == got.Token {
		t.Errorf("session = %+v, %v; want only the hash of token %q", got.Session, err, got.Token)
	}
}

func TestNewSession_WrapsRandomError(t *testing.T) {
	errEntropy := errors.New("no entropy")

	got, err := NewSession(failingRandom{errEntropy}, SessionOpening{UserID: 7, Now: testNow})

	if !errors.Is(err, errEntropy) || err.Error() == errEntropy.Error() || got != (OpenedSession{}) {
		t.Errorf("NewSession = %+v, %v; want no session and the error wrapped with context", got, err)
	}
}

func TestSession_IsExpiredAt(t *testing.T) {
	session := Session{ExpiresAt: testNow}

	if session.IsExpiredAt(testNow.Add(-time.Second)) || !session.IsExpiredAt(testNow) {
		t.Errorf("IsExpiredAt: want valid before %v and expired from it", testNow)
	}
}
