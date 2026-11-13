package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

func assertKeyedError(t *testing.T, err error, sentinel error, want shared.MessageKey, wantParams map[string]any) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want %v", err, sentinel)
	}
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) || keyed.Key != want {
		t.Fatalf("error = %v, want key %q", err, want)
	}
	if !reflect.DeepEqual(keyed.Params, wantParams) {
		t.Errorf("params = %v, want %v", keyed.Params, wantParams)
	}
}

func assertNothingParsedOrSaved(t *testing.T, fixture importFixture) {
	t.Helper()
	if len(fixture.bci.parsedTexts) != 0 || len(fixture.other.parsedTexts) != 0 {
		t.Errorf("parsed bci %v, other %v; want nothing parsed", fixture.bci.parsedTexts, fixture.other.parsedTexts)
	}
	if len(fixture.statements.created) != 0 {
		t.Errorf("Create called %d times, want 0", len(fixture.statements.created))
	}
}

func detectInput() ImportStatementInput {
	input := bciInput()
	input.Bank = ""
	return input
}

func TestImportStatement_UnregisteredChosenBankIsUnknownAndNotParsed(t *testing.T) {
	fixture := newImportFixture()
	input := bciInput()
	input.Bank = "santander"

	_, err := fixture.useCase.Execute(context.Background(), userOne, input)

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyBankUnknown, nil)
	assertNothingParsedOrSaved(t, fixture)
}

func TestImportStatement_WithoutBankUsesTheOnlyParserThatDetectsIt(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.detects = true

	got, err := fixture.useCase.Execute(context.Background(), userOne, detectInput())

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if got.Statement.Bank != "bci" {
		t.Errorf("statement bank = %q, want bci", got.Statement.Bank)
	}
	if !reflect.DeepEqual(fixture.bci.parsedTexts, []string{"statement text"}) || len(fixture.other.parsedTexts) != 0 {
		t.Errorf("bci parsed %v, other parsed %v; want only bci", fixture.bci.parsedTexts, fixture.other.parsedTexts)
	}
}

func TestImportStatement_WithoutBankNoParserDetectsIt(t *testing.T) {
	fixture := newImportFixture()

	_, err := fixture.useCase.Execute(context.Background(), userOne, detectInput())

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyBankNotDetected, map[string]any{"count": 0})
	assertNothingParsedOrSaved(t, fixture)
}

func TestImportStatement_WithoutBankTwoParsersDetectIt(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.detects = true
	fixture.other.detects = true

	_, err := fixture.useCase.Execute(context.Background(), userOne, detectInput())

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyBankNotDetected, map[string]any{"count": 2})
	assertNothingParsedOrSaved(t, fixture)
}

func TestImportStatement_ChosenRegisteredBankSkipsDetect(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.detects = true

	if _, err := fixture.useCase.Execute(context.Background(), userOne, bciInput()); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	if fixture.bci.detectCalls != 0 || fixture.other.detectCalls != 0 {
		t.Errorf("Detect called bci %d, other %d times; want 0", fixture.bci.detectCalls, fixture.other.detectCalls)
	}
}
