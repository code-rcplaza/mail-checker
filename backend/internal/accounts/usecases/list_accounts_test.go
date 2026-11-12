package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestListAccounts_ReturnsTheRepositoryAccounts(t *testing.T) {
	syncedAt := testNow.Add(-1)
	saved := []domain.Account{
		{ID: 1, PublicID: "acct-1", Email: "user@example.com", PasswordEnc: []byte("sealed:a"), CreatedAt: testNow},
		{ID: 2, PublicID: "acct-2", Email: "other@example.com", PasswordEnc: []byte("sealed:b"), CreatedAt: testNow, LastSyncAt: &syncedAt, LastErrorKey: shared.MessageKey("emailsync.mail_auth_failed")},
	}
	accounts := &fakeAccounts{listed: saved}

	got, err := NewListAccounts(accounts).Execute(context.Background())

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	assertSameAccounts(t, got, saved)
}

func TestListAccounts_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")
	accounts := &fakeAccounts{listErr: cause}

	_, err := NewListAccounts(accounts).Execute(context.Background())

	assertWrapped(t, err, cause, "list accounts")
}
