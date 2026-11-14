package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

var _ ports.ExpiredAccess = ExpiredAccess{}

func TestExpiredAccess_DeletesOnlyWhatExpired(t *testing.T) {
	db, anaID, betoID := twoUsers(t)
	ctx := context.Background()
	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	codes, sessions, sends := NewLoginCodes(db), NewSessions(db), NewCodeSends(db)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	must(codes.Save(ctx, domain.LoginCode{UserID: anaID, CodeHash: "expired", ExpiresAt: now}))
	must(codes.Save(ctx, domain.LoginCode{UserID: betoID, CodeHash: "valid", ExpiresAt: now.Add(time.Second)}))
	must(sessions.Create(ctx, domain.Session{UserID: anaID, TokenHash: "expired", ExpiresAt: now.Add(-time.Hour)}))
	must(sessions.Create(ctx, domain.Session{UserID: anaID, TokenHash: "valid", ExpiresAt: now.Add(time.Hour)}))
	must(sends.Record(ctx, domain.CodeSend{Email: "ana@example.com", SentAt: now.Add(-2 * time.Hour)}))
	must(sends.Record(ctx, domain.CodeSend{Email: "ana@example.com", SentAt: now.Add(-time.Hour)}))
	must(sends.Record(ctx, domain.CodeSend{Email: "ana@example.com", SentAt: now.Add(-59 * time.Minute)}))
	cleanup := NewExpiredAccess(db)

	for range 2 {
		if err := cleanup.DeleteExpired(ctx, domain.CleanupCutoffAt(now)); err != nil {
			t.Fatalf("DeleteExpired error = %v, want nil (nothing left to delete the second time)", err)
		}
	}

	if _, err := codes.Find(ctx, anaID); err == nil {
		t.Errorf("the expired code is still there")
	}
	if _, err := codes.Find(ctx, betoID); err != nil {
		t.Errorf("the valid code was deleted: %v", err)
	}
	if _, err := sessions.FindByTokenHash(ctx, "expired"); err == nil {
		t.Errorf("the expired session is still there")
	}
	if _, err := sessions.FindByTokenHash(ctx, "valid"); err != nil {
		t.Errorf("the valid session was deleted: %v", err)
	}
	var sendsLeft int
	if err := db.QueryRow(`SELECT COUNT(*) FROM code_sends`).Scan(&sendsLeft); err != nil || sendsLeft != 1 {
		t.Errorf("sends left = %d, %v; want only the one of the last hour", sendsLeft, err)
	}
}
