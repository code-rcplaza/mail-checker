package sqlite

import (
	"context"
	"database/sql"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type PriceAlertDecisions struct {
	queries *queries.Queries
}

func NewPriceAlertDecisions(db *sql.DB) PriceAlertDecisions {
	return PriceAlertDecisions{queries: queries.New(db)}
}

// List gives every kept amount of each line, oldest first.
func (p PriceAlertDecisions) List(ctx context.Context, userID int64) (map[domain.GroupKey][]shared.Money, error) {
	rows, err := p.queries.ListKeptRises(ctx, userID)
	if err != nil {
		return nil, wrapped("list kept rises", err)
	}
	kept := make(map[domain.GroupKey][]shared.Money)
	for _, row := range rows {
		key := domain.GroupKey(row.LineKey)
		kept[key] = append(kept[key], shared.Money(row.Amount))
	}
	return kept, nil
}

// Save adds the amount to the line's history; keeping it again changes nothing.
func (p PriceAlertDecisions) Save(ctx context.Context, userID int64, kept domain.KeptRise) error {
	return wrapped("keep rise", p.queries.KeepRise(ctx, queries.KeepRiseParams{UserID: userID, LineKey: string(kept.Key), Amount: int64(kept.NewAmount)}))
}

type ExchangeRates struct {
	queries *queries.Queries
}

func NewExchangeRates(db *sql.DB) ExchangeRates { return ExchangeRates{queries: queries.New(db)} }

// List gives the user's rates by date; the domain picks the latest.
func (r ExchangeRates) List(ctx context.Context, userID int64) ([]shared.ExchangeRate, error) {
	rows, err := r.queries.ListExchangeRates(ctx, userID)
	if err != nil {
		return nil, wrapped("list exchange rates", err)
	}
	rates := make([]shared.ExchangeRate, 0, len(rows))
	for _, row := range rows {
		rates = append(rates, shared.ExchangeRate{PesosPerDollar: shared.Money(row.PesosPerDollar), Date: sharedsqlite.TimeFromUnix(row.Date)})
	}
	return rates, nil
}
