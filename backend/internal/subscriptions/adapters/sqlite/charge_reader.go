// Package sqlite reads what subscriptions needs from the app database.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
)

// ChargeReader reads the charges table directly, so subscriptions does not import statements.
type ChargeReader struct {
	queries *queries.Queries
}

func NewChargeReader(db *sql.DB) ChargeReader {
	return ChargeReader{queries: queries.New(db)}
}

// ListSince includes the charges at exactly since.
func (r ChargeReader) ListSince(ctx context.Context, userID int64, since time.Time) ([]shared.Charge, error) {
	rows, err := r.queries.ListChargesSince(ctx, queries.ListChargesSinceParams{UserID: userID, OccurredAt: since.Unix()})
	if err != nil {
		return nil, fmt.Errorf("list charges: %w", err)
	}
	charges := make([]shared.Charge, 0, len(rows))
	for _, row := range rows {
		charges = append(charges, chargeFromRow(row))
	}
	return charges, nil
}

func chargeFromRow(row queries.Charge) shared.Charge {
	return shared.Charge{
		ID:          row.ID,
		UserID:      row.UserID,
		Bank:        shared.Bank(row.Bank),
		Source:      shared.Source(row.Source),
		ExternalID:  row.ExternalID,
		StatementID: optionalID(row.StatementID),
		OccurredAt:  sharedsqlite.TimeFromUnix(row.OccurredAt),
		Amount:      shared.Money(row.Amount),
		Currency:    shared.Currency(row.Currency),
		Merchant:    row.Merchant,
		Card:        shared.Card(row.Card),
		CardLast4:   row.CardLast4,
	}
}

func optionalID(id sql.NullInt64) *int64 {
	if !id.Valid {
		return nil
	}
	return &id.Int64
}
