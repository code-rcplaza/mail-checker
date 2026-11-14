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
	dependencies ImportStatementDependencies
}

// ImportStatementDependencies are the ports and settings an import needs.
type ImportStatementDependencies struct {
	Extractor   ports.FileExtractor
	Parsers     []ports.BankStatementParser
	Statements  ports.StatementRepository
	Rates       ports.ExchangeRateRecorder
	NewPublicID func() string
	Now         func() time.Time
	Location    *time.Location // the user's calendar, to tell which months were read
}

// ImportResult is the saved statement and the months its charges cover; the
// months are calculated, not stored.
type ImportResult struct {
	Statement domain.Statement
	Months    []domain.Month
}

// statementText is the text of a file and the parser of its bank.
type statementText struct {
	text   string
	parser ports.BankStatementParser
}

// readStatement is everything read from the uploaded file before parsing it.
type readStatement struct {
	statementText
	filename    string
	kind        domain.FileKind
	billingDate time.Time
}

// NewImportStatement copies the bank registry so it stays read-only after main builds it.
func NewImportStatement(dependencies ImportStatementDependencies) ImportStatement {
	dependencies.Parsers = slices.Clone(dependencies.Parsers)
	return ImportStatement{dependencies: dependencies}
}

func (u ImportStatement) Execute(ctx context.Context, userID int64, input ImportStatementInput) (ImportResult, error) {
	read, err := u.readFile(input)
	if err != nil {
		return ImportResult{}, err
	}
	imported := u.parseStatement(userID, read)
	if err := rejectWithoutCharges(*imported.Statement); err != nil {
		return ImportResult{}, err
	}
	if err := u.recordExchangeRate(ctx, userID, read); err != nil {
		return ImportResult{}, err
	}
	if err := u.saveStatement(ctx, userID, imported); err != nil {
		return ImportResult{}, err
	}
	return ImportResult{Statement: *imported.Statement, Months: domain.MonthsRead(imported.Charges, u.dependencies.Location)}, nil
}

// readFile validates the upload, extracts its text, picks the bank and the
// billing date: everything parsing needs.
func (u ImportStatement) readFile(input ImportStatementInput) (readStatement, error) {
	if err := validateInput(input); err != nil {
		return readStatement{}, err
	}
	text, kind, err := u.extractText(input)
	if err != nil {
		return readStatement{}, err
	}
	parser, err := u.chooseParser(input.Bank, text)
	if err != nil {
		return readStatement{}, err
	}
	source := statementText{text: text, parser: parser}
	billingDate, err := chooseBillingDate(source, input.BillingDate)
	if err != nil {
		return readStatement{}, err
	}
	return readStatement{statementText: source, filename: input.Filename, kind: kind, billingDate: billingDate}, nil
}

func validateInput(input ImportStatementInput) error {
	if len(input.Data) == 0 {
		return shared.Invalid(domain.KeyFileEmpty, nil)
	}
	return nil
}

// extractText also rejects files whose content has no text, like a scanned PDF.
func (u ImportStatement) extractText(input ImportStatementInput) (string, domain.FileKind, error) {
	text, kind, err := u.dependencies.Extractor.Extract(input.Filename, input.Data)
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
	for _, parser := range u.dependencies.Parsers {
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
	for _, parser := range u.dependencies.Parsers {
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
func chooseBillingDate(source statementText, userDate time.Time) (time.Time, error) {
	if fileDate, hasFileDate := source.parser.BillingDate(source.text); hasFileDate {
		return fileDate, nil
	}
	if userDate.IsZero() {
		return time.Time{}, shared.Invalid(domain.KeyBillingDateRequired, nil)
	}
	return userDate, nil
}

// parseStatement builds the statement and its charges; the public id is
// generated first because every charge's external id uses it.
func (u ImportStatement) parseStatement(userID int64, read readStatement) domain.ImportedStatement {
	source := domain.StatementSource{UserID: userID, Bank: read.parser.Bank(), PublicID: u.dependencies.NewPublicID(), BillingDate: read.billingDate}
	lines, unread := read.parser.Parse(read.text)
	charges, invalid := domain.ChargesFromLines(lines, source)
	statement := domain.Statement{
		PublicID:    source.PublicID,
		UserID:      userID,
		Bank:        source.Bank,
		Filename:    read.filename,
		Kind:        read.kind,
		BillingDate: read.billingDate,
		ImportedAt:  u.dependencies.Now(),
		Charges:     len(charges),
		BadLines:    domain.MergeBadLines(unread, invalid),
	}
	return domain.ImportedStatement{Statement: &statement, Charges: charges}
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
func (u ImportStatement) recordExchangeRate(ctx context.Context, userID int64, read readStatement) error {
	pair, hasPair := read.parser.TransferPair(read.text)
	if !hasPair {
		return nil
	}
	paid := shared.DollarsPaidInPesos{Dollars: pair.Dollars, Pesos: pair.Pesos}
	rate, hasRate, err := shared.DeriveExchangeRate(paid, read.billingDate)
	if err != nil || !hasRate {
		return err
	}
	if err := u.dependencies.Rates.Record(ctx, userID, rate); err != nil {
		return fmt.Errorf("record exchange rate: %w", err)
	}
	return nil
}

func (u ImportStatement) saveStatement(ctx context.Context, userID int64, imported domain.ImportedStatement) error {
	if err := u.dependencies.Statements.Create(ctx, userID, imported); err != nil {
		return fmt.Errorf("save statement: %w", err)
	}
	return nil
}
