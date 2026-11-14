package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/sqlite/queries"
)

type ExchangeRates struct {
	queries *queries.Queries
}

func NewExchangeRates(db *sql.DB) ExchangeRates { return ExchangeRates{queries: queries.New(db)} }

// Record keeps the first rate of each date: recording the date again changes nothing.
func (r ExchangeRates) Record(ctx context.Context, userID int64, rate shared.ExchangeRate) error {
	err := r.queries.RecordExchangeRate(ctx, queries.RecordExchangeRateParams{UserID: userID, PesosPerDollar: int64(rate.PesosPerDollar), Date: rate.Date.Unix()})
	if err != nil {
		return fmt.Errorf("record exchange rate: %w", err)
	}
	return nil
}
