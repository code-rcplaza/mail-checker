package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestRemoveAccount_DeletesThePublicID(t *testing.T) {
	accounts := &fakeAccounts{}

	if err := NewRemoveAccount(accounts).Execute(context.Background(), "acct-1"); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	if !reflect.DeepEqual(accounts.deleted, []string{"acct-1"}) {
		t.Errorf("deleted %v, want [acct-1]", accounts.deleted)
	}
}

func TestRemoveAccount_UnknownPublicIDIsNotFound(t *testing.T) {
	accounts := &fakeAccounts{deleteErr: shared.ErrNotFound}

	err := NewRemoveAccount(accounts).Execute(context.Background(), "acct-404")

	assertKeyedError(t, err, shared.ErrNotFound, domain.KeyNotFound)
}

func TestRemoveAccount_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")
	accounts := &fakeAccounts{deleteErr: cause}

	err := NewRemoveAccount(accounts).Execute(context.Background(), "acct-1")

	assertWrapped(t, err, cause, "delete account")
}
