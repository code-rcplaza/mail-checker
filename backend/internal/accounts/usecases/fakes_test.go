package usecases

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var testNow = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)

// fakeCipher seals by reversing and prefixing, so the sealed bytes never
// contain the plain text.
type fakeCipher struct {
	encrypted  []string
	encryptErr error
}

func (f *fakeCipher) Encrypt(plain string) ([]byte, error) {
	f.encrypted = append(f.encrypted, plain)
	if f.encryptErr != nil {
		return nil, f.encryptErr
	}
	runes := []rune(plain)
	slices.Reverse(runes)
	return []byte("sealed:" + string(runes)), nil
}

func (f *fakeCipher) Decrypt(sealed []byte) (string, error) {
	runes := []rune(strings.TrimPrefix(string(sealed), "sealed:"))
	slices.Reverse(runes)
	return string(runes), nil
}

// fakeAccounts assigns IDs from 1, like the real repository.
type fakeAccounts struct {
	created   []domain.Account
	createErr error
	listed    []domain.Account
	listErr   error
}

func (f *fakeAccounts) Create(_ context.Context, account *domain.Account) error {
	if f.createErr != nil {
		return f.createErr
	}
	account.ID = int64(len(f.created) + 1)
	f.created = append(f.created, *account)
	return nil
}

func (f *fakeAccounts) List(context.Context) ([]domain.Account, error) {
	return f.listed, f.listErr
}

func assertKeyedError(t *testing.T, err error, sentinel error, want shared.MessageKey) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want %v", err, sentinel)
	}
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) || keyed.Key != want {
		t.Fatalf("error = %v, want key %q", err, want)
	}
}

func assertWrapped(t *testing.T, err, cause error, context string) {
	t.Helper()
	if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), context+": ") {
		t.Errorf("error = %v, want %q wrapping %v", err, context, cause)
	}
}

func assertSameAccounts(t *testing.T, got, want []domain.Account) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("accounts = %+v, want %+v", got, want)
	}
}
