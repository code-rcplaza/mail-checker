package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
)

var (
	_ ports.LoginCodes = LoginCodes{}
	_ ports.CodeSends  = CodeSends{}
	_ ports.Sessions   = Sessions{}
)

var expiresAt = time.Date(2026, 11, 13, 12, 40, 0, 0, time.UTC)

// twoUsers opens a database with Ana and Beto and returns their ids.
func twoUsers(t *testing.T) (*sql.DB, int64, int64) {
	t.Helper()
	db := sqlitetest.OpenMigrated(t)
	users := NewUsers(db, func() time.Time { return createdAt })
	anaUser, errAna := users.Create(context.Background(), ana())
	beto := ana()
	beto.PublicID, beto.Email = "0193f1a2-0000-7000-8000-000000000009", "beto@example.com"
	betoUser, errBeto := users.Create(context.Background(), beto)
	if errAna != nil || errBeto != nil {
		t.Fatalf("create users: %v, %v", errAna, errBeto)
	}
	return db, anaUser.ID, betoUser.ID
}

func TestLoginCodes_SaveReplacesThePreviousCodeOfTheUser(t *testing.T) {
	db, anaID, betoID := twoUsers(t)
	codes := NewLoginCodes(db)
	ctx := context.Background()
	first := domain.LoginCode{UserID: anaID, CodeHash: "hash-1", ExpiresAt: expiresAt, Attempts: 3}
	second := domain.LoginCode{UserID: anaID, CodeHash: "hash-2", ExpiresAt: expiresAt.Add(time.Minute)}
	betos := domain.LoginCode{UserID: betoID, CodeHash: "hash-beto", ExpiresAt: expiresAt}
	for _, code := range []domain.LoginCode{first, betos, second} {
		if err := codes.Save(ctx, code); err != nil {
			t.Fatalf("Save error = %v, want nil", err)
		}
	}

	gotAna, errAna := codes.Find(ctx, anaID)
	gotBeto, errBeto := codes.Find(ctx, betoID)

	if errAna != nil || gotAna != second {
		t.Errorf("Ana's code = %+v, %v; want %+v", gotAna, errAna, second)
	}
	if errBeto != nil || gotBeto != betos {
		t.Errorf("Beto's code = %+v, %v; want %+v", gotBeto, errBeto, betos)
	}
}

func TestLoginCodes_AddFailedAttemptCountsUp(t *testing.T) {
	db, anaID, _ := twoUsers(t)
	codes := NewLoginCodes(db)
	ctx := context.Background()
	if err := codes.Save(ctx, domain.LoginCode{UserID: anaID, CodeHash: "hash", ExpiresAt: expiresAt}); err != nil {
		t.Fatalf("Save error = %v", err)
	}

	first, errFirst := codes.AddFailedAttempt(ctx, anaID)
	second, errSecond := codes.AddFailedAttempt(ctx, anaID)
	stored, errFind := codes.Find(ctx, anaID)

	if errFirst != nil || errSecond != nil || errFind != nil || first != 1 || second != 2 || stored.Attempts != 2 {
		t.Errorf("attempts = %d, %d, stored %d (%v, %v, %v); want 1, 2 and 2 stored", first, second, stored.Attempts, errFirst, errSecond, errFind)
	}
}

func TestLoginCodes_MissingCodeIsNotFound(t *testing.T) {
	db, anaID, _ := twoUsers(t)
	codes := NewLoginCodes(db)
	ctx := context.Background()
	if err := codes.Save(ctx, domain.LoginCode{UserID: anaID, CodeHash: "hash", ExpiresAt: expiresAt}); err != nil {
		t.Fatalf("Save error = %v", err)
	}

	for range 2 {
		if err := codes.Delete(ctx, anaID); err != nil {
			t.Fatalf("Delete error = %v, want nil (idempotent)", err)
		}
	}
	_, errFind := codes.Find(ctx, anaID)
	_, errAttempt := codes.AddFailedAttempt(ctx, anaID)

	if !errors.Is(errFind, shared.ErrNotFound) || !errors.Is(errAttempt, shared.ErrNotFound) {
		t.Errorf("after delete: Find = %v, AddFailedAttempt = %v; want ErrNotFound for both", errFind, errAttempt)
	}
}

func TestCodeSends_ListSinceCountsOnlyTheRecentSendsOfThatEmail(t *testing.T) {
	sends := NewCodeSends(sqlitetest.OpenMigrated(t))
	ctx := context.Background()
	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	recent := []time.Time{now.Add(-50 * time.Minute), now.Add(-20 * time.Minute), now.Add(-5 * time.Minute)}
	records := []domain.CodeSend{
		{Email: "ana@example.com", SentAt: now.Add(-2 * time.Hour)},
		{Email: "ana@example.com", SentAt: recent[0]},
		{Email: "beto@example.com", SentAt: now.Add(-10 * time.Minute)},
		{Email: "ana@example.com", SentAt: recent[1]},
		{Email: "ana@example.com", SentAt: recent[2]},
	}
	for _, send := range records {
		if err := sends.Record(ctx, send); err != nil {
			t.Fatalf("Record error = %v, want nil", err)
		}
	}

	got, err := sends.ListSince(ctx, "ana@example.com", now.Add(-time.Hour))

	if err != nil || !reflect.DeepEqual(got, recent) {
		t.Errorf("ListSince = %v, %v; want the 3 sends of the last hour %v", got, err, recent)
	}
}

func TestSessions_CreateFindAndDelete(t *testing.T) {
	db, anaID, _ := twoUsers(t)
	sessions := NewSessions(db)
	ctx := context.Background()
	session := domain.Session{UserID: anaID, TokenHash: "token-hash", ExpiresAt: expiresAt}
	if err := sessions.Create(ctx, session); err != nil {
		t.Fatalf("Create error = %v, want nil", err)
	}

	got, err := sessions.FindByTokenHash(ctx, "token-hash")
	if err != nil || got != session {
		t.Fatalf("FindByTokenHash = %+v, %v; want %+v", got, err, session)
	}
	for range 2 {
		if err := sessions.Delete(ctx, "token-hash"); err != nil {
			t.Fatalf("Delete error = %v, want nil (idempotent)", err)
		}
	}
	_, err = sessions.FindByTokenHash(ctx, "token-hash")
	if !errors.Is(err, shared.ErrNotFound) {
		t.Errorf("FindByTokenHash after delete error = %v, want ErrNotFound", err)
	}
}

func TestSessions_UnknownHashIsNotFound(t *testing.T) {
	_, err := NewSessions(sqlitetest.OpenMigrated(t)).FindByTokenHash(context.Background(), "never-seen")

	if !errors.Is(err, shared.ErrNotFound) {
		t.Errorf("FindByTokenHash error = %v, want ErrNotFound", err)
	}
}
