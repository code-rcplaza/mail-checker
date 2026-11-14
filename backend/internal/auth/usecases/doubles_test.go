package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var testNow = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

var errDatabase = errors.New("database is locked")

// testClock lets a test move time forward between steps.
type testClock struct{ now time.Time }

func (c *testClock) Now() time.Time { return c.now }

func (c *testClock) advance(by time.Duration) { c.now = c.now.Add(by) }

// fixedRandom fills every read with the same byte: codes and tokens are predictable.
type fixedRandom byte

func (f fixedRandom) Read(p []byte) (int, error) {
	copy(p, bytes.Repeat([]byte{byte(f)}, len(p)))
	return len(p), nil
}

var _ ports.UserRepository = (*usersByEmail)(nil)

type usersByEmail struct {
	users map[string]domain.User
	calls int
}

func usersWith(users ...domain.User) *usersByEmail {
	byEmail := make(map[string]domain.User)
	for _, user := range users {
		byEmail[user.Email] = user
	}
	return &usersByEmail{users: byEmail}
}

func (u *usersByEmail) Create(_ context.Context, user domain.User) (domain.User, error) {
	u.calls++
	user.ID = int64(len(u.users) + 1)
	u.users[user.Email] = user
	return user, nil
}

func (u *usersByEmail) FindByEmail(_ context.Context, email string) (domain.User, error) {
	u.calls++
	user, found := u.users[email]
	if !found {
		return domain.User{}, fmt.Errorf("find user: %w", shared.ErrNotFound)
	}
	return user, nil
}

func (u *usersByEmail) FindByID(_ context.Context, id int64) (domain.User, error) {
	u.calls++
	for _, user := range u.users {
		if user.ID == id {
			return user, nil
		}
	}
	return domain.User{}, fmt.Errorf("find user: %w", shared.ErrNotFound)
}

var _ ports.LoginCodes = (*codesByUser)(nil)

type codesByUser struct {
	codes   map[int64]domain.LoginCode
	saveErr error
}

func newCodesByUser() *codesByUser { return &codesByUser{codes: make(map[int64]domain.LoginCode)} }

func (c *codesByUser) Save(_ context.Context, code domain.LoginCode) error {
	if c.saveErr != nil {
		return c.saveErr
	}
	c.codes[code.UserID] = code
	return nil
}

func (c *codesByUser) Find(_ context.Context, userID int64) (domain.LoginCode, error) {
	code, found := c.codes[userID]
	if !found {
		return domain.LoginCode{}, fmt.Errorf("find code: %w", shared.ErrNotFound)
	}
	return code, nil
}

func (c *codesByUser) AddFailedAttempt(_ context.Context, userID int64) (int, error) {
	code := c.codes[userID]
	code.Attempts++
	c.codes[userID] = code
	return code.Attempts, nil
}

func (c *codesByUser) Delete(_ context.Context, userID int64) error {
	delete(c.codes, userID)
	return nil
}

var _ ports.CodeSends = sendsByEmail{}

type sendsByEmail map[string][]time.Time

func (s sendsByEmail) ListSince(_ context.Context, email string, since time.Time) ([]time.Time, error) {
	var after []time.Time
	for _, sentAt := range s[email] {
		if sentAt.After(since) {
			after = append(after, sentAt)
		}
	}
	return after, nil
}

func (s sendsByEmail) Record(_ context.Context, send domain.CodeSend) error {
	s[send.Email] = append(s[send.Email], send.SentAt)
	return nil
}

var _ ports.CodeSender = (*recordingSender)(nil)

type sentCode struct{ email, code string }

type recordingSender struct {
	sent []sentCode
	err  error
}

func (r *recordingSender) Send(_ context.Context, email, code string) error {
	if r.err != nil {
		return r.err
	}
	r.sent = append(r.sent, sentCode{email, code})
	return nil
}

var _ ports.Sessions = sessionsByHash{}

type sessionsByHash map[string]domain.Session

func (s sessionsByHash) Create(_ context.Context, session domain.Session) error {
	s[session.TokenHash] = session
	return nil
}

func (s sessionsByHash) FindByTokenHash(_ context.Context, tokenHash string) (domain.Session, error) {
	session, found := s[tokenHash]
	if !found {
		return domain.Session{}, fmt.Errorf("find session: %w", shared.ErrNotFound)
	}
	return session, nil
}

func (s sessionsByHash) Delete(_ context.Context, tokenHash string) error {
	delete(s, tokenHash)
	return nil
}

var ana = domain.User{ID: 1, PublicID: "user-ana", Name: "Ana", Email: "ana@example.com", AcceptedTermsAt: testNow}

// deliveryFixture is a CodeDelivery with every double at hand.
type deliveryFixture struct {
	codes  *codesByUser
	sends  sendsByEmail
	sender *recordingSender
	clock  *testClock
}

func newDeliveryFixture() deliveryFixture {
	return deliveryFixture{codes: newCodesByUser(), sends: sendsByEmail{}, sender: &recordingSender{}, clock: &testClock{now: testNow}}
}

func (f deliveryFixture) delivery() CodeDelivery {
	return NewCodeDelivery(CodeDeliveryDependencies{Codes: f.codes, Sends: f.sends, Sender: f.sender, Random: fixedRandom(1), Now: f.clock.Now})
}

func assertKeyedError(t *testing.T, err error, kind error, key shared.MessageKey) {
	t.Helper()
	var keyed shared.KeyedError
	if !errors.Is(err, kind) || !errors.As(err, &keyed) || keyed.Key != key {
		t.Errorf("error = %v, want %v with key %q", err, kind, key)
	}
}

func assertWrapped(t *testing.T, err, cause error) {
	t.Helper()
	if !errors.Is(err, cause) || err.Error() == cause.Error() {
		t.Errorf("error = %v, want %v wrapped with context", err, cause)
	}
}
