package sqlite

import (
	"context"
	"database/sql"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type Manuals struct {
	queries *queries.Queries
}

func NewManuals(db *sql.DB) Manuals { return Manuals{queries: queries.New(db)} }

func (m Manuals) List(ctx context.Context, userID int64) ([]domain.ManualSubscription, error) {
	rows, err := m.queries.ListManuals(ctx, userID)
	if err != nil {
		return nil, wrapped("list manual subscriptions", err)
	}
	manuals := make([]domain.ManualSubscription, 0, len(rows))
	for _, row := range rows {
		manuals = append(manuals, domain.ManualSubscription{
			PublicID: row.PublicID, Name: row.Name, Amount: shared.Money(row.Amount), Currency: shared.Currency(row.Currency), BillingDay: int(row.BillingDay),
		})
	}
	return manuals, nil
}

func (m Manuals) Create(ctx context.Context, userID int64, manual domain.ManualSubscription) error {
	return wrapped("create manual subscription", m.queries.CreateManual(ctx, queries.CreateManualParams{
		UserID: userID, PublicID: manual.PublicID, Name: manual.Name, Amount: int64(manual.Amount), Currency: string(manual.Currency), BillingDay: int64(manual.BillingDay),
	}))
}

func (m Manuals) Update(ctx context.Context, userID int64, manual domain.ManualSubscription) error {
	rows, err := m.queries.UpdateManual(ctx, queries.UpdateManualParams{
		Name: manual.Name, Amount: int64(manual.Amount), Currency: string(manual.Currency), BillingDay: int64(manual.BillingDay), UserID: userID, PublicID: manual.PublicID,
	})
	if err != nil {
		return wrapped("update manual subscription", err)
	}
	return notFoundWhenNone("update manual subscription", rows)
}

func (m Manuals) Delete(ctx context.Context, userID int64, publicID string) error {
	rows, err := m.queries.DeleteManual(ctx, queries.DeleteManualParams{UserID: userID, PublicID: publicID})
	if err != nil {
		return wrapped("delete manual subscription", err)
	}
	return notFoundWhenNone("delete manual subscription", rows)
}
