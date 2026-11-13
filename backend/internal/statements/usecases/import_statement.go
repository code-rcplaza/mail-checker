package usecases

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

type ImportStatementInput struct {
	Filename    string
	Data        []byte
	BillingDate time.Time   // optional: only used when the file shows none
	Bank        shared.Bank // empty: detect it from the text
}

// ImportStatement turns an uploaded statement into charges, following the
// numbered steps of ARCHITECTURE.md.
type ImportStatement struct {
	extractor   ports.FileExtractor
	parsers     []ports.BankStatementParser
	statements  ports.StatementRepository
	rates       ports.ExchangeRateRecorder
	newPublicID func() string
	now         func() time.Time
	location    *time.Location // the user's calendar, to tell which months were read
}

// ImportResult is the saved statement and the months its charges cover; the
// months are calculated, not stored.
type ImportResult struct {
	Statement domain.Statement
	Months    []domain.Month
}

// NewImportStatement copies the bank registry so it stays read-only after main builds it.
func NewImportStatement(
	extractor ports.FileExtractor,
	parsers []ports.BankStatementParser,
	statements ports.StatementRepository,
	rates ports.ExchangeRateRecorder,
	newPublicID func() string,
	now func() time.Time,
	location *time.Location,
) ImportStatement {
	return ImportStatement{
		extractor:   extractor,
		parsers:     slices.Clone(parsers),
		statements:  statements,
		rates:       rates,
		newPublicID: newPublicID,
		now:         now,
		location:    location,
	}
}

func (u ImportStatement) Execute(ctx context.Context, userID int64, input ImportStatementInput) (ImportResult, error) {
	if err := validateInput(input); err != nil {
		return ImportResult{}, err
	}
	text, kind, err := u.extractText(input)
	if err != nil {
		return ImportResult{}, err
	}
	parser, err := u.chooseParser(input.Bank, text)
	if err != nil {
		return ImportResult{}, err
	}
	billingDate, err := chooseBillingDate(parser, text, input.BillingDate)
	if err != nil {
		return ImportResult{}, err
	}
	statement, charges := u.parseStatement(userID, input.Filename, kind, parser, text, billingDate)
	if err := rejectWithoutCharges(statement); err != nil {
		return ImportResult{}, err
	}
	if err := u.recordExchangeRate(ctx, userID, parser, text, billingDate); err != nil {
		return ImportResult{}, err
	}
	if err := u.saveStatement(ctx, userID, &statement, charges); err != nil {
		return ImportResult{}, err
	}
	return ImportResult{Statement: statement, Months: domain.MonthsRead(charges, u.location)}, nil
}

func validateInput(input ImportStatementInput) error {
	if len(input.Data) == 0 {
		return shared.Invalid(domain.KeyFileEmpty, nil)
	}
	return nil
}

// extractText also rejects files whose content has no text, like a scanned PDF.
func (u ImportStatement) extractText(input ImportStatementInput) (string, domain.FileKind, error) {
	text, kind, err := u.extractor.Extract(input.Filename, input.Data)
	if err != nil {
		return "", "", fmt.Errorf("extract statement text: %w", err)
	}
	if strings.TrimSpace(text) == "" {
		return "", "", shared.Invalid(domain.KeyFileEmpty, nil)
	}
	return text, kind, nil
}

// chooseParser uses the chosen bank, or detects it when none was chosen.
func (u ImportStatement) chooseParser(bank shared.Bank, text string) (ports.BankStatementParser, error) {
	if bank == "" {
		return u.detectParser(text)
	}
	return u.findParser(bank)
}

func (u ImportStatement) findParser(bank shared.Bank) (ports.BankStatementParser, error) {
	for _, parser := range u.parsers {
		if parser.Bank() == bank {
			return parser, nil
		}
	}
	return nil, shared.Invalid(domain.KeyBankUnknown, nil)
}

// detectParser accepts only an unambiguous match: guessing between two banks
// would import charges under the wrong one.
func (u ImportStatement) detectParser(text string) (ports.BankStatementParser, error) {
	var detected []ports.BankStatementParser
	for _, parser := range u.parsers {
		if parser.Detect(text) {
			detected = append(detected, parser)
		}
	}
	if len(detected) != 1 {
		return nil, shared.Invalid(domain.KeyBankNotDetected, map[string]any{"count": len(detected)})
	}
	return detected[0], nil
}

// chooseBillingDate prefers the date the file shows: the user typed theirs
// only in case the file had none.
func chooseBillingDate(parser ports.BankStatementParser, text string, userDate time.Time) (time.Time, error) {
	if fileDate, hasFileDate := parser.BillingDate(text); hasFileDate {
		return fileDate, nil
	}
	if userDate.IsZero() {
		return time.Time{}, shared.Invalid(domain.KeyBillingDateRequired, nil)
	}
	return userDate, nil
}

// parseStatement builds the statement and its charges; the public id is
// generated first because every charge's external id uses it.
func (u ImportStatement) parseStatement(userID int64, filename string, kind domain.FileKind, parser ports.BankStatementParser, text string, billingDate time.Time) (domain.Statement, []shared.Charge) {
	publicID := u.newPublicID()
	lines, unread := parser.Parse(text)
	charges, invalid := domain.ChargesFromLines(lines, userID, parser.Bank(), publicID, billingDate)
	statement := domain.Statement{
		PublicID:    publicID,
		UserID:      userID,
		Bank:        parser.Bank(),
		Filename:    filename,
		Kind:        kind,
		BillingDate: billingDate,
		ImportedAt:  u.now(),
		Charges:     len(charges),
		BadLines:    domain.MergeBadLines(unread, invalid),
	}
	return statement, charges
}

// rejectWithoutCharges keeps a statement with no charges out of the base, so a
// misread file does not block a correct reimport of the same billing date.
func rejectWithoutCharges(statement domain.Statement) error {
	if statement.Charges > 0 {
		return nil
	}
	return shared.Invalid(domain.KeyNoCharges, map[string]any{"badLines": len(statement.BadLines)})
}

// recordExchangeRate keeps the rate of the transfer pair, dated on the
// billing date. It runs before the statement is saved: if it fails nothing is
// imported, and if saving fails afterwards the rate stays, which is harmless
// since a rate per date is recorded only once.
func (u ImportStatement) recordExchangeRate(ctx context.Context, userID int64, parser ports.BankStatementParser, text string, billingDate time.Time) error {
	pair, hasPair := parser.TransferPair(text)
	if !hasPair {
		return nil
	}
	rate, hasRate, err := shared.DeriveExchangeRate(pair.Pesos, pair.Dollars, billingDate)
	if err != nil || !hasRate {
		return err
	}
	if err := u.rates.Record(ctx, userID, rate); err != nil {
		return fmt.Errorf("record exchange rate: %w", err)
	}
	return nil
}

func (u ImportStatement) saveStatement(ctx context.Context, userID int64, statement *domain.Statement, charges []shared.Charge) error {
	if err := u.statements.Create(ctx, userID, statement, charges); err != nil {
		return fmt.Errorf("save statement: %w", err)
	}
	return nil
}
