package usecases

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

var testNow = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)

type fakeDirectory struct {
	credentials []domain.Credential
}

func (f fakeDirectory) Credentials(context.Context) ([]domain.Credential, error) {
	return f.credentials, nil
}

// fakeMailReader serves emails by account and sender, and records every query.
type fakeMailReader struct {
	emails  map[string][]domain.RawEmail // key: account email + "|" + sender
	queries []domain.MailQuery
}

func (f *fakeMailReader) Read(_ context.Context, credential domain.Credential, query domain.MailQuery) ([]domain.RawEmail, error) {
	f.queries = append(f.queries, query)
	return f.emails[credential.Email+"|"+query.From], nil
}

// fakeParser recognizes the emails from its bank's sender and takes the
// merchant from the subject.
type fakeParser struct {
	bank   shared.Bank
	sender string
}

func (f fakeParser) Bank() shared.Bank { return f.bank }

func (f fakeParser) Query() domain.MailQuery {
	return domain.MailQuery{From: f.sender, Subject: "compra"}
}

func (f fakeParser) Parse(email domain.RawEmail) (shared.ParsedCharge, bool) {
	if email.From != f.sender {
		return shared.ParsedCharge{}, false
	}
	return shared.ParsedCharge{Amount: 599000, Currency: shared.CurrencyCLP, Merchant: email.Subject, Card: "credito"}, true
}

// fakeChargeWriter is idempotent by external id, like the real Insert.
type fakeChargeWriter struct {
	inserted []shared.Charge
}

func (f *fakeChargeWriter) Insert(_ context.Context, charge shared.Charge) (bool, error) {
	for _, existing := range f.inserted {
		if existing.Source == charge.Source && existing.ExternalID == charge.ExternalID {
			return false, nil
		}
	}
	f.inserted = append(f.inserted, charge)
	return true, nil
}

type fakeUnparsed struct {
	saved     map[string]domain.UnparsedEmail
	forgotten []string
}

func newFakeUnparsed() *fakeUnparsed {
	return &fakeUnparsed{saved: make(map[string]domain.UnparsedEmail)}
}

func (f *fakeUnparsed) Save(_ context.Context, email domain.UnparsedEmail) error {
	f.saved[email.ExternalID] = email
	return nil
}

func (f *fakeUnparsed) Forget(_ context.Context, externalID string) error {
	f.forgotten = append(f.forgotten, externalID)
	delete(f.saved, externalID)
	return nil
}

var (
	bankA = fakeParser{bank: "bank-a", sender: "avisos@bank-a.example"}
	bankB = fakeParser{bank: "bank-b", sender: "avisos@bank-b.example"}
)

func credential(id int64, email string) domain.Credential {
	return domain.Credential{AccountID: id, AccountPublicID: "account-" + email, Email: email, Password: "app-password"}
}

func purchase(messageID, sender, merchant string) domain.RawEmail {
	return domain.RawEmail{MessageID: messageID, Subject: merchant, From: sender, ReceivedAt: testNow.AddDate(0, 0, -5), Body: "compra"}
}

type syncFixture struct {
	reader   *fakeMailReader
	writer   *fakeChargeWriter
	unparsed *fakeUnparsed
	useCase  SyncAll
}

func newSyncFixture(credentials []domain.Credential, emails map[string][]domain.RawEmail, parsers ...ports.BankEmailParser) syncFixture {
	reader := &fakeMailReader{emails: emails}
	writer := &fakeChargeWriter{}
	unparsed := newFakeUnparsed()
	useCase := NewSyncAll(fakeDirectory{credentials: credentials}, reader, parsers, writer, unparsed, func() time.Time { return testNow })
	return syncFixture{reader: reader, writer: writer, unparsed: unparsed, useCase: useCase}
}

func TestSyncAll_ReadsWithBankFilterSinceNowMinusMonths(t *testing.T) {
	fixture := newSyncFixture([]domain.Credential{credential(1, "one@example.com")}, nil, bankA)

	if _, err := fixture.useCase.Execute(context.Background(), 3); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	want := []domain.MailQuery{{From: bankA.sender, Subject: "compra", Since: time.Date(2026, 8, 12, 12, 0, 0, 0, time.UTC)}}
	if !reflect.DeepEqual(fixture.reader.queries, want) {
		t.Errorf("queries = %+v, want %+v", fixture.reader.queries, want)
	}
}

func TestSyncAll_InsertsRecognizedEmailsAndCountsThem(t *testing.T) {
	emails := map[string][]domain.RawEmail{
		"one@example.com|" + bankA.sender: {
			purchase("<a1@bank-a>", bankA.sender, "STREAMING CO"),
			purchase("<a2@bank-a>", bankA.sender, "MUSIC APP"),
		},
	}
	fixture := newSyncFixture([]domain.Credential{credential(1, "one@example.com")}, emails, bankA)

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if len(fixture.writer.inserted) != 2 {
		t.Fatalf("inserted = %d charges, want 2", len(fixture.writer.inserted))
	}
	for _, c := range fixture.writer.inserted {
		if c.Source != shared.SourceEmail {
			t.Errorf("charge %s source = %q, want %q", c.ExternalID, c.Source, shared.SourceEmail)
		}
	}
	if len(results) != 1 || results[0].Found != 2 || results[0].New != 2 {
		t.Errorf("results = %+v, want one result with Found 2 and New 2", results)
	}
}

func TestSyncAll_OneResultPerAccountInDirectoryOrder(t *testing.T) {
	credentials := []domain.Credential{credential(1, "one@example.com"), credential(2, "two@example.com")}
	emails := map[string][]domain.RawEmail{
		"two@example.com|" + bankA.sender: {purchase("<a1@bank-a>", bankA.sender, "STREAMING CO")},
	}
	fixture := newSyncFixture(credentials, emails, bankA)

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	want := []domain.SyncResult{
		{AccountPublicID: "account-one@example.com", Email: "one@example.com"},
		{AccountPublicID: "account-two@example.com", Email: "two@example.com", Found: 1, New: 1},
	}
	if !reflect.DeepEqual(results, want) {
		t.Errorf("results = %+v, want %+v", results, want)
	}
}

func TestSyncAll_SearchesEachBankAndParsesWithItsParser(t *testing.T) {
	emails := map[string][]domain.RawEmail{
		"one@example.com|" + bankA.sender: {purchase("<a1@bank-a>", bankA.sender, "STREAMING CO")},
		"one@example.com|" + bankB.sender: {purchase("<b1@bank-b>", bankB.sender, "MUSIC APP")},
	}
	fixture := newSyncFixture([]domain.Credential{credential(1, "one@example.com")}, emails, bankA, bankB)

	if _, err := fixture.useCase.Execute(context.Background(), 3); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	senders := make([]string, 0, len(fixture.reader.queries))
	for _, q := range fixture.reader.queries {
		senders = append(senders, q.From)
	}
	if !reflect.DeepEqual(senders, []string{bankA.sender, bankB.sender}) {
		t.Errorf("searched senders = %v, want both banks", senders)
	}
	banksByID := make(map[string]shared.Bank)
	for _, c := range fixture.writer.inserted {
		banksByID[c.ExternalID] = c.Bank
	}
	want := map[string]shared.Bank{"<a1@bank-a>": "bank-a", "<b1@bank-b>": "bank-b"}
	if !reflect.DeepEqual(banksByID, want) {
		t.Errorf("charge banks = %v, want %v", banksByID, want)
	}
}
