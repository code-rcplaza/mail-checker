package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const plainPassword = "abcdefghijklmnop"

type addFixture struct {
	accounts *fakeAccounts
	cipher   *fakeCipher
	useCase  AddAccount
}

func newAddFixture() addFixture {
	accounts := &fakeAccounts{}
	cipher := &fakeCipher{}
	useCase := NewAddAccount(accounts, cipher, func() string { return "acct-1" }, func() time.Time { return testNow })
	return addFixture{accounts: accounts, cipher: cipher, useCase: useCase}
}

func validAddInput() AddAccountInput {
	return AddAccountInput{Email: "  User@Example.com ", AppPassword: "abcd efgh ijkl mnop"}
}

func TestAddAccount_EncryptsNormalizedPasswordAndCreatesOnce(t *testing.T) {
	fixture := newAddFixture()

	if _, err := fixture.useCase.Execute(context.Background(), validAddInput()); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	if !reflect.DeepEqual(fixture.cipher.encrypted, []string{plainPassword}) {
		t.Errorf("encrypted %q, want [%s]", fixture.cipher.encrypted, plainPassword)
	}
	sealed, _ := (&fakeCipher{}).Encrypt(plainPassword)
	want := []domain.Account{{ID: 1, PublicID: "acct-1", Email: "user@example.com", PasswordEnc: sealed, CreatedAt: testNow}}
	assertSameAccounts(t, fixture.accounts.created, want)
}

func TestAddAccount_PlainPasswordIsNotInSavedOrReturnedAccount(t *testing.T) {
	fixture := newAddFixture()

	returned, err := fixture.useCase.Execute(context.Background(), validAddInput())

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if len(fixture.accounts.created) != 1 {
		t.Fatalf("Create called %d times, want 1", len(fixture.accounts.created))
	}
	for name, account := range map[string]domain.Account{"saved": fixture.accounts.created[0], "returned": returned} {
		if len(account.PasswordEnc) == 0 {
			t.Errorf("%s account has no sealed password", name)
		}
		if bytes.Contains(account.PasswordEnc, []byte(plainPassword)) || strings.Contains(fmt.Sprintf("%+v", account), plainPassword) {
			t.Errorf("%s account %+v contains the plain password", name, account)
		}
	}
}

func TestAddAccount_InvalidInputIsRejectedWithoutEncryptingOrSaving(t *testing.T) {
	tests := []struct {
		name  string
		input AddAccountInput
		want  shared.MessageKey
	}{
		{"email", AddAccountInput{Email: "user@", AppPassword: "abcd efgh ijkl mnop"}, domain.KeyEmailInvalid},
		{"app password", AddAccountInput{Email: "user@example.com", AppPassword: "abcd1234"}, domain.KeyAppPasswordInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newAddFixture()

			_, err := fixture.useCase.Execute(context.Background(), tt.input)

			assertKeyedError(t, err, shared.ErrInvalid, tt.want)
			if len(fixture.cipher.encrypted) != 0 || len(fixture.accounts.created) != 0 {
				t.Errorf("encrypted %d, created %d; want nothing", len(fixture.cipher.encrypted), len(fixture.accounts.created))
			}
		})
	}
}

func TestAddAccount_EmailTakenConflictKeepsKindAndKey(t *testing.T) {
	fixture := newAddFixture()
	fixture.accounts.createErr = shared.Conflict(domain.KeyEmailTaken)

	_, err := fixture.useCase.Execute(context.Background(), validAddInput())

	assertKeyedError(t, err, shared.ErrConflict, domain.KeyEmailTaken)
}

func TestAddAccount_CipherAndRepositoryErrorsAreWrapped(t *testing.T) {
	cause := errors.New("boom")
	tests := []struct {
		name        string
		breakIt     func(addFixture)
		wantContext string
	}{
		{"cipher", func(f addFixture) { f.cipher.encryptErr = cause }, "encrypt app password"},
		{"repository", func(f addFixture) { f.accounts.createErr = cause }, "save account"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newAddFixture()
			tt.breakIt(fixture)

			_, err := fixture.useCase.Execute(context.Background(), validAddInput())

			assertWrapped(t, err, cause, tt.wantContext)
		})
	}
}
