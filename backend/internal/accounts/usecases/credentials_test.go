package usecases

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func sealedAccount(id int64, publicID, email, password string) domain.Account {
	sealed, _ := (&fakeCipher{}).Encrypt(password)
	return domain.Account{ID: id, PublicID: publicID, Email: email, PasswordEnc: sealed, CreatedAt: testNow}
}

func TestCredentials_OnePerAccountInRepositoryOrder(t *testing.T) {
	accounts := &fakeAccounts{listed: []domain.Account{
		sealedAccount(2, "acct-2", "other@example.com", "qrstuvwxyzabcdef"),
		sealedAccount(1, "acct-1", "user@example.com", plainPassword),
	}}

	got, err := NewCredentials(accounts, &fakeCipher{}).Execute(context.Background())

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	want := []domain.Credential{
		{AccountID: 2, AccountPublicID: "acct-2", Email: "other@example.com", Password: "qrstuvwxyzabcdef"},
		{AccountID: 1, AccountPublicID: "acct-1", Email: "user@example.com", Password: plainPassword},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("credentials = %+v, want %+v", got, want)
	}
}

func TestCredentials_UnreadablePasswordIsConflictWithoutPartialResult(t *testing.T) {
	broken := sealedAccount(2, "acct-2", "other@example.com", "")
	broken.PasswordEnc = []byte("sealed with a lost key")
	accounts := &fakeAccounts{listed: []domain.Account{
		sealedAccount(1, "acct-1", "user@example.com", plainPassword),
		broken,
	}}

	got, err := NewCredentials(accounts, &fakeCipher{}).Execute(context.Background())

	if got != nil {
		t.Errorf("credentials = %+v, want nil", got)
	}
	assertKeyedError(t, err, shared.ErrConflict, domain.KeyPasswordUnreadable)
	if strings.Contains(err.Error(), plainPassword) {
		t.Errorf("error %q contains a password", err)
	}
}

func TestCredentials_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")
	accounts := &fakeAccounts{listErr: cause}

	_, err := NewCredentials(accounts, &fakeCipher{}).Execute(context.Background())

	assertWrapped(t, err, cause, "list accounts")
}
