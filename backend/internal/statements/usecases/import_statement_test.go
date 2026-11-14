package usecases

import (
	"context"
	"reflect"
	"slices"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

var (
	testNow     = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	billingDate = time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
	chargeDate  = time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)
)

type fakeExtractor struct {
	text        string
	kind        domain.FileKind
	err         error
	gotFilename string
	gotData     []byte
	calls       int
}

func (f *fakeExtractor) Extract(filename string, data []byte) (string, domain.FileKind, error) {
	f.calls++
	f.gotFilename, f.gotData = filename, data
	return f.text, f.kind, f.err
}

// fakeStatementParser returns fixed lines and records the text it parsed.
type fakeStatementParser struct {
	bank         shared.Bank
	detects      bool
	lines        []domain.StatementLine
	bad          []int
	parsedTexts  []string
	detectCalls  int
	billingDate  time.Time            // zero: the file shows no billing date
	transferPair *domain.TransferPair // nil: the file shows no pair
}

func (f *fakeStatementParser) Bank() shared.Bank { return f.bank }

func (f *fakeStatementParser) Detect(string) bool {
	f.detectCalls++
	return f.detects
}

func (f *fakeStatementParser) Parse(text string) ([]domain.StatementLine, []int) {
	f.parsedTexts = append(f.parsedTexts, text)
	return f.lines, f.bad
}

func (f *fakeStatementParser) BillingDate(string) (time.Time, bool) {
	return f.billingDate, !f.billingDate.IsZero()
}

func (f *fakeStatementParser) TransferPair(string) (domain.TransferPair, bool) {
	if f.transferPair == nil {
		return domain.TransferPair{}, false
	}
	return *f.transferPair, true
}

type createdStatement struct {
	statement domain.Statement
	charges   []shared.Charge
}

// fakeStatements assigns IDs from 1, like the real repository.
type fakeStatements struct {
	created   []createdStatement
	createErr error
	listed    []domain.Statement
	listErr   error
	deleted   []string
	deleteErr error
}

func (f *fakeStatements) Create(_ context.Context, userID int64, imported domain.ImportedStatement) error {
	if f.createErr != nil {
		return f.createErr
	}
	if f.hasStatement(userID, *imported.Statement) {
		return shared.Conflict(domain.KeyAlreadyImported)
	}
	imported.Statement.ID = int64(len(f.created) + 1)
	f.created = append(f.created, createdStatement{statement: *imported.Statement, charges: imported.Charges})
	return nil
}

// hasStatement mirrors the real uniqueness: user, bank and billing date.
func (f *fakeStatements) hasStatement(userID int64, statement domain.Statement) bool {
	return slices.ContainsFunc(f.created, func(c createdStatement) bool {
		return c.statement.UserID == userID && c.statement.Bank == statement.Bank && c.statement.BillingDate.Equal(statement.BillingDate)
	})
}

// List returns the fixed list when a test sets one; otherwise what the user created.
func (f *fakeStatements) List(_ context.Context, userID int64) ([]domain.Statement, error) {
	if f.listed != nil || f.listErr != nil {
		return f.listed, f.listErr
	}
	var mine []domain.Statement
	for _, c := range f.created {
		if c.statement.UserID == userID {
			mine = append(mine, c.statement)
		}
	}
	return mine, nil
}

func (f *fakeStatements) Delete(_ context.Context, userID int64, publicID string) error {
	f.deleted = append(f.deleted, publicID)
	if f.deleteErr != nil {
		return f.deleteErr
	}
	if f.belongsToAnotherUser(userID, publicID) {
		return shared.ErrNotFound
	}
	return nil
}

func (f *fakeStatements) belongsToAnotherUser(userID int64, publicID string) bool {
	return slices.ContainsFunc(f.created, func(c createdStatement) bool {
		return c.statement.PublicID == publicID && c.statement.UserID != userID
	})
}

func statementLine(number int, merchant string) domain.StatementLine {
	return domain.StatementLine{
		Number: number,
		Charge: shared.ParsedCharge{Amount: 599000, Currency: shared.CurrencyCLP, Merchant: merchant, Card: "credito", Date: chargeDate},
	}
}

type importFixture struct {
	extractor  *fakeExtractor
	bci        *fakeStatementParser
	other      *fakeStatementParser
	statements *fakeStatements
	rates      *fakeExchangeRates
	useCase    ImportStatement
}

// newImportFixture registers two banks; bci has two valid lines (2 and 4),
// one invalid line (3) and one it could not read (5).
func newImportFixture() importFixture {
	extractor := &fakeExtractor{text: "statement text", kind: domain.FileKindPDF}
	bci := &fakeStatementParser{
		bank:  "bci",
		lines: []domain.StatementLine{statementLine(2, "STREAMING CO"), statementLine(3, ""), statementLine(4, "MUSIC APP")},
		bad:   []int{5},
	}
	other := &fakeStatementParser{bank: "other-bank"}
	statements := &fakeStatements{}
	rates := &fakeExchangeRates{}
	useCase := NewImportStatement(ImportStatementDependencies{
		Extractor: extractor, Parsers: []ports.BankStatementParser{other, bci}, Statements: statements, Rates: rates,
		NewPublicID: func() string { return "stmt-1" }, Now: func() time.Time { return testNow }, Location: time.UTC,
	})
	return importFixture{extractor: extractor, bci: bci, other: other, statements: statements, rates: rates, useCase: useCase}
}

func bciInput() ImportStatementInput {
	return ImportStatementInput{Filename: "cartola-octubre.pdf", Data: []byte("%PDF-fake"), BillingDate: billingDate, Bank: "bci"}
}

func TestImportStatement_ExtractsParsesWithChosenBankAndCreatesOnce(t *testing.T) {
	fixture := newImportFixture()

	if _, err := fixture.useCase.Execute(context.Background(), userOne, bciInput()); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	if fixture.extractor.gotFilename != "cartola-octubre.pdf" || string(fixture.extractor.gotData) != "%PDF-fake" {
		t.Errorf("extracted %q (%q), want the uploaded file", fixture.extractor.gotFilename, fixture.extractor.gotData)
	}
	if !reflect.DeepEqual(fixture.bci.parsedTexts, []string{"statement text"}) || len(fixture.other.parsedTexts) != 0 {
		t.Errorf("bci parsed %v, other parsed %v; want only bci with the extracted text", fixture.bci.parsedTexts, fixture.other.parsedTexts)
	}
	if len(fixture.statements.created) != 1 {
		t.Fatalf("Create called %d times, want 1", len(fixture.statements.created))
	}
	ids := make([]string, 0, 2)
	for _, c := range fixture.statements.created[0].charges {
		ids = append(ids, c.ExternalID)
	}
	if !reflect.DeepEqual(ids, []string{"stmt-1:2", "stmt-1:4"}) {
		t.Errorf("created charges = %v, want [stmt-1:2 stmt-1:4]", ids)
	}
}

func TestImportStatement_SavedStatementHasAllItsData(t *testing.T) {
	fixture := newImportFixture()

	if _, err := fixture.useCase.Execute(context.Background(), userOne, bciInput()); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	if len(fixture.statements.created) != 1 {
		t.Fatalf("Create called %d times, want 1", len(fixture.statements.created))
	}
	want := domain.Statement{
		ID:          1,
		PublicID:    "stmt-1",
		UserID:      userOne,
		Bank:        "bci",
		Filename:    "cartola-octubre.pdf",
		Kind:        domain.FileKindPDF,
		BillingDate: billingDate,
		ImportedAt:  testNow,
		Charges:     2,
		BadLines:    []int{3, 5},
	}
	if got := fixture.statements.created[0].statement; !reflect.DeepEqual(got, want) {
		t.Errorf("saved statement = %+v, want %+v", got, want)
	}
}

func TestImportStatement_ReturnsTheSavedStatement(t *testing.T) {
	fixture := newImportFixture()

	got, err := fixture.useCase.Execute(context.Background(), userOne, bciInput())

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if len(fixture.statements.created) != 1 {
		t.Fatalf("Create called %d times, want 1", len(fixture.statements.created))
	}
	if saved := fixture.statements.created[0].statement; !reflect.DeepEqual(got.Statement, saved) {
		t.Errorf("returned %+v, want the saved %+v", got, saved)
	}
}
