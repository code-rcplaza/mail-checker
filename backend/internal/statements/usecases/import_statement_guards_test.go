package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

func TestImportStatement_EmptyFileIsRejectedBeforeExtracting(t *testing.T) {
	fixture := newImportFixture()
	input := bciInput()
	input.Data = nil

	_, err := fixture.useCase.Execute(context.Background(), input)

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyFileEmpty, nil)
	if fixture.extractor.calls != 0 {
		t.Errorf("Extract called %d times, want 0", fixture.extractor.calls)
	}
	assertNothingParsedOrSaved(t, fixture)
}

func TestImportStatement_BlankExtractedTextIsEmptyFile(t *testing.T) {
	for _, text := range []string{"", " \n\t "} {
		fixture := newImportFixture()
		fixture.extractor.text = text

		_, err := fixture.useCase.Execute(context.Background(), bciInput())

		assertKeyedError(t, err, shared.ErrInvalid, domain.KeyFileEmpty, nil)
		assertNothingParsedOrSaved(t, fixture)
	}
}

func TestImportStatement_ZeroBillingDateIsMissing(t *testing.T) {
	fixture := newImportFixture()
	input := bciInput()
	input.BillingDate = time.Time{}

	_, err := fixture.useCase.Execute(context.Background(), input)

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyBillingDateMissing, nil)
	assertNothingParsedOrSaved(t, fixture)
}

func TestImportStatement_NoValidChargeIsRejectedWithBadLineCount(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.lines = []domain.StatementLine{statementLine(3, "")}
	fixture.bci.bad = []int{5}

	_, err := fixture.useCase.Execute(context.Background(), bciInput())

	assertKeyedError(t, err, shared.ErrInvalid, domain.KeyNoCharges, map[string]any{"badLines": 2})
	if len(fixture.statements.created) != 0 {
		t.Errorf("Create called %d times, want 0", len(fixture.statements.created))
	}
}

func TestImportStatement_AlreadyImportedConflictKeepsKindAndKey(t *testing.T) {
	fixture := newImportFixture()
	fixture.statements.createErr = shared.Conflict(domain.KeyAlreadyImported)

	_, err := fixture.useCase.Execute(context.Background(), bciInput())

	assertKeyedError(t, err, shared.ErrConflict, domain.KeyAlreadyImported, nil)
}

func TestImportStatement_ExtractorAndRepositoryErrorsAreWrapped(t *testing.T) {
	cause := errors.New("boom")
	tests := []struct {
		name        string
		breakIt     func(importFixture)
		wantContext string
	}{
		{"extractor", func(f importFixture) { f.extractor.err = cause }, "extract statement text"},
		{"repository", func(f importFixture) { f.statements.createErr = cause }, "save statement"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newImportFixture()
			tt.breakIt(fixture)

			_, err := fixture.useCase.Execute(context.Background(), bciInput())

			if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), tt.wantContext+": ") {
				t.Errorf("error = %v, want %q wrapping %v", err, tt.wantContext, cause)
			}
		})
	}
}
