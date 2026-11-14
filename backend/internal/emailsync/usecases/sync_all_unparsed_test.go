package usecases

import (
	"context"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
)

const accountOne, accountTwo = "one@example.com", "two@example.com"

func TestSyncAll_SameEmailInTwoAccountsIsInsertedOnce(t *testing.T) {
	email := purchase("<a1@bank-a>", bankA.sender, "STREAMING CO")
	emails := map[string][]domain.RawEmail{
		accountOne + "|" + bankA.sender: {email},
		accountTwo + "|" + bankA.sender: {email},
	}
	fixture := newSyncFixture([]domain.Credential{credential(1, accountOne), credential(2, accountTwo)}, emails, bankA)

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if results[0].Found != 1 || results[0].New != 1 || results[1].Found != 1 || results[1].New != 0 {
		t.Errorf("results = %+v, want Found 1 in both and New 1 then 0", results)
	}
	if len(fixture.writer.inserted) != 1 {
		t.Errorf("inserted = %d charges, want 1", len(fixture.writer.inserted))
	}
}

func TestSyncAll_EmailSavedInAPreviousSyncIsNotNew(t *testing.T) {
	emails := map[string][]domain.RawEmail{accountOne + "|" + bankA.sender: {purchase("<a1@bank-a>", bankA.sender, "STREAMING CO")}}
	fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, emails, bankA)
	if _, err := fixture.useCase.Execute(context.Background(), 3); err != nil {
		t.Fatalf("first Execute error = %v, want nil", err)
	}

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("second Execute error = %v, want nil", err)
	}
	if results[0].New != 0 || results[0].Found != 1 {
		t.Errorf("second sync result = %+v, want Found 1 and New 0", results[0])
	}
}

func TestSyncAll_UnrecognizedEmailIsSavedAsUnparsed(t *testing.T) {
	// Arrives through bank A's filter but does not look like bank A's notice.
	stranger := purchase("<x1@other>", "newsletter@other.example", "OFFERS")
	emails := map[string][]domain.RawEmail{accountOne + "|" + bankA.sender: {stranger}}
	fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, emails, bankA)

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if _, isSaved := fixture.unparsed.saved["<x1@other>"]; !isSaved {
		t.Errorf("unparsed saved = %v, want <x1@other>", fixture.unparsed.saved)
	}
	if results[0].Unparsed != 1 || results[0].New != 0 {
		t.Errorf("result = %+v, want Unparsed 1 and New 0", results[0])
	}
}

func TestSyncAll_PreviouslyUnparsedEmailNowRecognizedIsInsertedAndForgotten(t *testing.T) {
	email := purchase("<a1@bank-a>", bankA.sender, "STREAMING CO")
	emails := map[string][]domain.RawEmail{accountOne + "|" + bankA.sender: {email}}
	fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, emails, bankA)
	fixture.unparsed.saved[email.MessageID] = domain.UnparsedFromEmail(domain.BankEmail{Email: email, AccountID: 1, Bank: "bank-a"})

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if results[0].New != 1 {
		t.Errorf("result = %+v, want New 1", results[0])
	}
	if len(fixture.unparsed.forgotten) != 1 || fixture.unparsed.forgotten[0] != email.MessageID {
		t.Errorf("forgotten = %v, want [%s]", fixture.unparsed.forgotten, email.MessageID)
	}
}

func TestSyncAll_RecognizedEmailWithInvalidChargeIsSavedAsUnparsed(t *testing.T) {
	emptyMerchant := purchase("<a1@bank-a>", bankA.sender, "") // the fake parser takes the merchant from the subject
	emails := map[string][]domain.RawEmail{accountOne + "|" + bankA.sender: {emptyMerchant}}
	fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, emails, bankA)

	results, err := fixture.useCase.Execute(context.Background(), 3)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if _, isSaved := fixture.unparsed.saved["<a1@bank-a>"]; !isSaved {
		t.Errorf("unparsed saved = %v, want <a1@bank-a>", fixture.unparsed.saved)
	}
	if results[0].Unparsed != 1 || len(fixture.writer.inserted) != 0 {
		t.Errorf("result = %+v with %d inserted, want Unparsed 1 and nothing inserted", results[0], len(fixture.writer.inserted))
	}
}
