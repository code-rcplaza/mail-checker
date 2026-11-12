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
	BillingDate time.Time
	Bank        shared.Bank // empty: detect it from the text
}

// ImportStatement turns an uploaded statement into charges, following the
// numbered steps of ARCHITECTURE.md.
type ImportStatement struct {
	extractor   ports.FileExtractor
	parsers     []ports.BankStatementParser
	statements  ports.StatementRepository
	newPublicID func() string
	now         func() time.Time
}

// NewImportStatement copies the bank registry so it stays read-only after main builds it.
func NewImportStatement(
	extractor ports.FileExtractor,
	parsers []ports.BankStatementParser,
	statements ports.StatementRepository,
	newPublicID func() string,
	now func() time.Time,
) ImportStatement {
	return ImportStatement{
		extractor:   extractor,
		parsers:     slices.Clone(parsers),
		statements:  statements,
		newPublicID: newPublicID,
		now:         now,
	}
}

func (u ImportStatement) Execute(ctx context.Context, input ImportStatementInput) (domain.Statement, error) {
	if err := validateInput(input); err != nil {
		return domain.Statement{}, err
	}
	text, kind, err := u.extractText(input)
	if err != nil {
		return domain.Statement{}, err
	}
	parser, err := u.chooseParser(input.Bank, text)
	if err != nil {
		return domain.Statement{}, err
	}
	statement, charges := u.parseStatement(input, kind, parser, text)
	if err := rejectWithoutCharges(statement); err != nil {
		return domain.Statement{}, err
	}
	if err := u.saveStatement(ctx, &statement, charges); err != nil {
		return domain.Statement{}, err
	}
	return statement, nil
}

func validateInput(input ImportStatementInput) error {
	if len(input.Data) == 0 {
		return shared.Invalid(domain.KeyFileEmpty, nil)
	}
	if input.BillingDate.IsZero() {
		return shared.Invalid(domain.KeyBillingDateMissing, nil)
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

// parseStatement builds the statement and its charges; the public id is
// generated first because every charge's external id uses it.
func (u ImportStatement) parseStatement(input ImportStatementInput, kind domain.FileKind, parser ports.BankStatementParser, text string) (domain.Statement, []shared.Charge) {
	publicID := u.newPublicID()
	lines, unread := parser.Parse(text)
	charges, invalid := domain.ChargesFromLines(lines, parser.Bank(), publicID, input.BillingDate)
	statement := domain.Statement{
		PublicID:    publicID,
		Bank:        parser.Bank(),
		Filename:    input.Filename,
		Kind:        kind,
		BillingDate: input.BillingDate,
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

func (u ImportStatement) saveStatement(ctx context.Context, statement *domain.Statement, charges []shared.Charge) error {
	if err := u.statements.Create(ctx, statement, charges); err != nil {
		return fmt.Errorf("save statement: %w", err)
	}
	return nil
}
