// Package sqlite stores the statements feature in the app database.
package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

type Statements struct {
	db      *sql.DB
	queries *queries.Queries
}

func NewStatements(db *sql.DB) Statements {
	return Statements{db: db, queries: queries.New(db)}
}

// Create saves the statement and its charges in one transaction: one rejected
// charge leaves nothing behind. The ID is assigned only once it is committed.
func (s Statements) Create(ctx context.Context, userID int64, imported domain.ImportedStatement) error {
	var statementID int64
	err := sharedsqlite.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		q := s.queries.WithTx(tx)
		id, err := createStatement(ctx, q, ownedStatement{userID: userID, statement: imported.Statement})
		if err != nil {
			return err
		}
		statementID = id
		return createCharges(ctx, q, statementCharges{statementID: id, charges: imported.Charges})
	})
	if err != nil {
		return err
	}
	imported.Statement.ID = statementID
	return nil
}

// ownedStatement takes the owner from the port, never from the statement itself.
type ownedStatement struct {
	userID    int64
	statement *domain.Statement
}

func createStatement(ctx context.Context, q *queries.Queries, owned ownedStatement) (int64, error) {
	statement := owned.statement
	badLines, err := json.Marshal(nonNilLines(statement.BadLines))
	if err != nil {
		return 0, fmt.Errorf("encode bad lines: %w", err)
	}
	id, err := q.CreateStatement(ctx, queries.CreateStatementParams{
		PublicID:    statement.PublicID,
		UserID:      owned.userID,
		Bank:        string(statement.Bank),
		Filename:    statement.Filename,
		Kind:        string(statement.Kind),
		BillingDate: statement.BillingDate.Unix(),
		ImportedAt:  statement.ImportedAt.Unix(),
		BadLines:    string(badLines),
	})
	if sharedsqlite.IsUniqueViolation(err) {
		return 0, shared.Conflict(domain.KeyAlreadyImported)
	}
	if err != nil {
		return 0, fmt.Errorf("create statement: %w", err)
	}
	return id, nil
}

func nonNilLines(lines []int) []int {
	if lines == nil {
		return []int{}
	}
	return lines
}

type statementCharges struct {
	statementID int64
	charges     []shared.Charge
}

func createCharges(ctx context.Context, q *queries.Queries, owned statementCharges) error {
	for _, charge := range owned.charges {
		err := q.CreateCharge(ctx, queries.CreateChargeParams{
			UserID:      charge.UserID,
			Bank:        string(charge.Bank),
			Source:      string(charge.Source),
			ExternalID:  charge.ExternalID,
			StatementID: sql.NullInt64{Int64: owned.statementID, Valid: true},
			OccurredAt:  charge.OccurredAt.Unix(),
			Amount:      int64(charge.Amount),
			Currency:    string(charge.Currency),
			Merchant:    charge.Merchant,
			Card:        string(charge.Card),
			CardLast4:   charge.CardLast4,
		})
		if err != nil {
			return fmt.Errorf("create charge %s: %w", charge.ExternalID, err)
		}
	}
	return nil
}

func (s Statements) List(ctx context.Context, userID int64) ([]domain.Statement, error) {
	rows, err := s.queries.ListStatements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list statements: %w", err)
	}
	statements := make([]domain.Statement, 0, len(rows))
	for _, row := range rows {
		statement, err := statementFromRow(row)
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func statementFromRow(row queries.ListStatementsRow) (domain.Statement, error) {
	var badLines []int
	if err := json.Unmarshal([]byte(row.BadLines), &badLines); err != nil {
		return domain.Statement{}, fmt.Errorf("decode bad lines of statement %s: %w", row.PublicID, err)
	}
	if len(badLines) == 0 {
		badLines = nil
	}
	return domain.Statement{
		ID:          row.ID,
		PublicID:    row.PublicID,
		UserID:      row.UserID,
		Bank:        shared.Bank(row.Bank),
		Filename:    row.Filename,
		Kind:        domain.FileKind(row.Kind),
		BillingDate: sharedsqlite.TimeFromUnix(row.BillingDate),
		ImportedAt:  sharedsqlite.TimeFromUnix(row.ImportedAt),
		Charges:     int(row.ChargeCount),
		BadLines:    badLines,
	}, nil
}

// Delete answers ErrNotFound without a key: the use case adds statements.not_found.
func (s Statements) Delete(ctx context.Context, userID int64, publicID string) error {
	deleted, err := s.queries.DeleteStatement(ctx, queries.DeleteStatementParams{UserID: userID, PublicID: publicID})
	if err != nil {
		return fmt.Errorf("delete statement: %w", err)
	}
	if deleted == 0 {
		return fmt.Errorf("delete statement: %w", shared.ErrNotFound)
	}
	return nil
}
