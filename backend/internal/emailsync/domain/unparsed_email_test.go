package domain

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestUnparsedFromEmail_KeepsEmailValues(t *testing.T) {
	email := purchaseEmail()

	got := UnparsedFromEmail(BankEmail{Email: email, AccountID: 7, Bank: "bci"})

	want := UnparsedEmail{
		ExternalID: email.MessageID,
		AccountID:  7,
		Bank:       "bci",
		Subject:    email.Subject,
		ReceivedAt: email.ReceivedAt,
		Preview:    email.Body,
	}
	if got != want {
		t.Errorf("UnparsedFromEmail = %+v, want %+v", got, want)
	}
}

func TestUnparsedFromEmail_LongBodyPreviewKeepsWholeCharacters(t *testing.T) {
	email := purchaseEmail()
	email.Body = strings.Repeat("ñá", 150) // 300 characters, 2 bytes each

	got := UnparsedFromEmail(BankEmail{Email: email, AccountID: 7, Bank: "bci"}).Preview

	if count := utf8.RuneCountInString(got); count != PreviewMaxRunes {
		t.Errorf("preview has %d characters, want %d", count, PreviewMaxRunes)
	}
	if !utf8.ValidString(got) {
		t.Errorf("preview splits a character: %q", got)
	}
	if !strings.HasPrefix(email.Body, got) {
		t.Errorf("preview is not the start of the body")
	}
}

func TestUnparsedFromEmail_ShortBodyIsWholePreview(t *testing.T) {
	email := purchaseEmail()

	got := UnparsedFromEmail(BankEmail{Email: email, AccountID: 7, Bank: "bci"}).Preview

	if got != email.Body {
		t.Errorf("preview = %q, want the whole body %q", got, email.Body)
	}
}
