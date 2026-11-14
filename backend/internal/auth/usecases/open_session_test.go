package usecases

import (
	"context"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func openSession(fixture verifyFixture, sessions sessionsByHash, random fixedRandom) (string, error) {
	useCase := NewOpenSession(OpenSessionDependencies{Verify: fixture.verifyCode(), Sessions: sessions, Random: random, Now: fixture.clock.Now})
	return useCase.Execute(context.Background(), CodeEntry{Email: ana.Email, Code: anasCode})
}

func TestOpenSession_RightCodeReturnsTheTokenAndStoresTheSession(t *testing.T) {
	sessions := sessionsByHash{}

	token, err := openSession(newVerifyFixture(), sessions, fixedRandom(7))

	want := domain.Session{UserID: ana.ID, TokenHash: domain.HashSecret(token), ExpiresAt: testNow.Add(domain.SessionLifetime)}
	if err != nil || token == "" || len(sessions) != 1 || sessions[want.TokenHash] != want {
		t.Errorf("token = %q, sessions = %+v, err = %v; want one session %+v", token, sessions, err, want)
	}
}

func TestOpenSession_WrongCodeOpensNoSession(t *testing.T) {
	fixture := newVerifyFixture()
	fixture.codes.codes[ana.ID] = domain.LoginCode{UserID: ana.ID, CodeHash: domain.HashSecret("999999"), ExpiresAt: testNow.Add(domain.CodeLifetime)}
	sessions := sessionsByHash{}

	token, err := openSession(fixture, sessions, fixedRandom(7))

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyCodeInvalid)
	if token != "" || len(sessions) != 0 {
		t.Errorf("token = %q, sessions = %+v; want none", token, sessions)
	}
}

func TestOpenSession_AnotherDeviceKeepsItsSession(t *testing.T) {
	sessions := sessionsByHash{}

	first, _ := openSession(newVerifyFixture(), sessions, fixedRandom(7))
	second, err := openSession(newVerifyFixture(), sessions, fixedRandom(8))

	if err != nil || first == second || len(sessions) != 2 {
		t.Errorf("tokens %q and %q, sessions = %d, err = %v; want two distinct sessions", first, second, len(sessions), err)
	}
}
