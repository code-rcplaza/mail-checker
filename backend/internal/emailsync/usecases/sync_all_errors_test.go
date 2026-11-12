package usecases

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
)

// twoAccountsWithOneEmailEach: each account has one recognizable bank A email.
func twoAccountsWithOneEmailEach() syncFixture {
	emails := map[string][]domain.RawEmail{
		accountOne + "|" + bankA.sender: {purchase("<a1@bank-a>", bankA.sender, "STREAMING CO")},
		accountTwo + "|" + bankA.sender: {purchase("<a2@bank-a>", bankA.sender, "MUSIC APP")},
	}
	return newSyncFixture([]domain.Credential{credential(1, accountOne), credential(2, accountTwo)}, emails, bankA)
}

func TestSyncAll_MailReadErrorStaysInItsAccount(t *testing.T) {
	cases := []struct {
		name    string
		readErr error
		wantKey string
	}{
		{"rejected credentials", fmt.Errorf("imap login: %w", domain.ErrMailAuthRejected), "emailsync.mail_auth_failed"},
		{"other read error", errors.New("connection reset"), "emailsync.mail_read_failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := twoAccountsWithOneEmailEach()
			fixture.reader.errs[accountOne] = tc.readErr

			results, err := fixture.useCase.Execute(context.Background(), 3)

			if err != nil {
				t.Fatalf("Execute error = %v, want nil (the error belongs to the account)", err)
			}
			if len(results) != 2 {
				t.Fatalf("results = %+v, want two", results)
			}
			if string(results[0].ErrorKey) != tc.wantKey {
				t.Errorf("first account ErrorKey = %q, want %q", results[0].ErrorKey, tc.wantKey)
			}
			if results[1].ErrorKey != "" || results[1].New != 1 {
				t.Errorf("second account = %+v, want synced with New 1", results[1])
			}
		})
	}
}

func TestSyncAll_RecordsEachAccountSyncWithClockTimeAndErrorKey(t *testing.T) {
	fixture := twoAccountsWithOneEmailEach()
	fixture.reader.errs[accountTwo] = errors.New("connection reset")

	if _, err := fixture.useCase.Execute(context.Background(), 3); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	want := []syncRecord{
		{accountID: 1, at: testNow, errKey: ""},
		{accountID: 2, at: testNow, errKey: domain.KeyMailReadFailed},
	}
	if !reflect.DeepEqual(fixture.directory.records, want) {
		t.Errorf("records = %+v, want %+v", fixture.directory.records, want)
	}
}

func TestSyncAll_StorageErrorsStopTheSyncWrapped(t *testing.T) {
	errDatabase := errors.New("database is locked")
	cases := []struct {
		name   string
		inject func(f syncFixture)
	}{
		{"credentials", func(f syncFixture) { f.directory.credentialsErr = errDatabase }},
		{"insert", func(f syncFixture) { f.writer.insertErr = errDatabase }},
		{"save unparsed", func(f syncFixture) {
			f.reader.emails[accountOne+"|"+bankA.sender] = []domain.RawEmail{purchase("<x1@other>", "newsletter@other.example", "OFFERS")}
			f.unparsed.saveErr = errDatabase
		}},
		{"forget unparsed", func(f syncFixture) { f.unparsed.forgetErr = errDatabase }},
		{"record sync", func(f syncFixture) { f.directory.recordErr = errDatabase }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := twoAccountsWithOneEmailEach()
			tc.inject(fixture)

			_, err := fixture.useCase.Execute(context.Background(), 3)

			if !errors.Is(err, errDatabase) {
				t.Fatalf("errors.Is(err, errDatabase) = false (err = %v), want true", err)
			}
			if err.Error() == errDatabase.Error() {
				t.Errorf("err = %q, want it wrapped with context", err)
			}
		})
	}
}
