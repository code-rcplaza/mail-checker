package accounts

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	accountsdomain "github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	accountsusecases "github.com/code-rcplaza/mail-checker/backend/internal/accounts/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var testNow = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)

// The bridge must satisfy the port emailsync declares.
var _ ports.AccountDirectory = AccountDirectory{}

// prefixCipher "seals" by prefixing; anything without the prefix is unreadable.
type prefixCipher struct{}

func (prefixCipher) Encrypt(plain string) ([]byte, error) { return []byte("sealed:" + plain), nil }

func (prefixCipher) Decrypt(sealed []byte) (string, error) {
	plain, found := strings.CutPrefix(string(sealed), "sealed:")
	if !found {
		return "", errors.New("message authentication failed")
	}
	return plain, nil
}

type syncMark struct {
	id     int64
	at     time.Time
	errKey shared.MessageKey
}

type fakeAccounts struct {
	listed []accountsdomain.Account
	marked []syncMark
}

func (f *fakeAccounts) Create(context.Context, *accountsdomain.Account) error { return nil }

func (f *fakeAccounts) List(context.Context) ([]accountsdomain.Account, error) { return f.listed, nil }

func (f *fakeAccounts) Delete(context.Context, string) error { return nil }

func (f *fakeAccounts) MarkSynced(_ context.Context, record accountsdomain.SyncRecord) error {
	f.marked = append(f.marked, syncMark{id: record.AccountID, at: record.At, errKey: record.ErrorKey})
	return nil
}

func newDirectory(accounts *fakeAccounts) AccountDirectory {
	return NewAccountDirectory(
		accountsusecases.NewCredentials(accounts, prefixCipher{}),
		accountsusecases.NewRecordSync(accounts),
	)
}

func account(id int64, publicID, email string, sealed string) accountsdomain.Account {
	return accountsdomain.Account{ID: id, PublicID: publicID, Email: email, PasswordEnc: []byte(sealed), CreatedAt: testNow}
}

func TestAccountDirectory_CredentialsAreEmailsyncCredentialsInOrder(t *testing.T) {
	accounts := &fakeAccounts{listed: []accountsdomain.Account{
		account(2, "acct-2", "other@example.com", "sealed:qrstuvwxyzabcdef"),
		account(1, "acct-1", "user@example.com", "sealed:abcdefghijklmnop"),
	}}

	got, err := newDirectory(accounts).Credentials(context.Background())

	if err != nil {
		t.Fatalf("Credentials error = %v, want nil", err)
	}
	want := []domain.Credential{
		{AccountID: 2, AccountPublicID: "acct-2", Email: "other@example.com", Password: "qrstuvwxyzabcdef"},
		{AccountID: 1, AccountPublicID: "acct-1", Email: "user@example.com", Password: "abcdefghijklmnop"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("credentials = %+v, want %+v", got, want)
	}
}

func TestAccountDirectory_CredentialsErrorKeepsKindAndKey(t *testing.T) {
	accounts := &fakeAccounts{listed: []accountsdomain.Account{
		account(1, "acct-1", "user@example.com", "sealed with a lost key"),
	}}

	_, err := newDirectory(accounts).Credentials(context.Background())

	var keyed shared.KeyedError
	if !errors.Is(err, shared.ErrConflict) || !errors.As(err, &keyed) || keyed.Key != accountsdomain.KeyPasswordUnreadable {
		t.Errorf("error = %v, want conflict with key %q", err, accountsdomain.KeyPasswordUnreadable)
	}
}

func TestAccountDirectory_RecordSyncForwardsItsValues(t *testing.T) {
	accounts := &fakeAccounts{}

	if err := newDirectory(accounts).RecordSync(context.Background(), domain.SyncRecord{AccountID: 7, At: testNow, ErrorKey: "emailsync.mail_auth_failed"}); err != nil {
		t.Fatalf("RecordSync error = %v, want nil", err)
	}

	want := []syncMark{{id: 7, at: testNow, errKey: "emailsync.mail_auth_failed"}}
	if !reflect.DeepEqual(accounts.marked, want) {
		t.Errorf("marked %+v, want %+v", accounts.marked, want)
	}
}
