package sqlite

import (
	"context"
	"database/sql"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type Overrides struct {
	queries *queries.Queries
}

func NewOverrides(db *sql.DB) Overrides { return Overrides{queries: queries.New(db)} }

func (o Overrides) List(ctx context.Context, userID int64) (map[domain.GroupKey]domain.LineOverride, error) {
	rows, err := o.queries.ListOverrides(ctx, userID)
	if err != nil {
		return nil, wrapped("list overrides", err)
	}
	overrides := make(map[domain.GroupKey]domain.LineOverride, len(rows))
	for _, row := range rows {
		overrides[domain.GroupKey(row.LineKey)] = domain.LineOverride{Amount: shared.Money(row.Amount.Int64), BillingDay: int(row.BillingDay.Int64)}
	}
	return overrides, nil
}

func (o Overrides) Save(ctx context.Context, userID int64, overridden domain.OverriddenLine) error {
	return wrapped("save override", o.queries.SaveOverride(ctx, queries.SaveOverrideParams{
		UserID:     userID,
		LineKey:    string(overridden.Key),
		Amount:     nullWhenZero(int64(overridden.Override.Amount)),
		BillingDay: nullWhenZero(int64(overridden.Override.BillingDay)),
	}))
}

func (o Overrides) Delete(ctx context.Context, userID int64, key domain.GroupKey) error {
	return wrapped("delete override", o.queries.DeleteOverride(ctx, queries.DeleteOverrideParams{UserID: userID, LineKey: string(key)}))
}
