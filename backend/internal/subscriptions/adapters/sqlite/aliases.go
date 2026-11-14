package sqlite

import (
	"context"
	"database/sql"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type Aliases struct {
	queries *queries.Queries
}

func NewAliases(db *sql.DB) Aliases { return Aliases{queries: queries.New(db)} }

func (a Aliases) List(ctx context.Context, userID int64) (map[domain.GroupKey]string, error) {
	rows, err := a.queries.ListAliases(ctx, userID)
	if err != nil {
		return nil, wrapped("list aliases", err)
	}
	aliases := make(map[domain.GroupKey]string, len(rows))
	for _, row := range rows {
		aliases[domain.GroupKey(row.LineKey)] = row.Alias
	}
	return aliases, nil
}

func (a Aliases) Save(ctx context.Context, userID int64, alias domain.LineAlias) error {
	return wrapped("save alias", a.queries.SaveAlias(ctx, queries.SaveAliasParams{UserID: userID, LineKey: string(alias.Key), Alias: alias.Alias}))
}

func (a Aliases) Delete(ctx context.Context, userID int64, key domain.GroupKey) error {
	return wrapped("delete alias", a.queries.DeleteAlias(ctx, queries.DeleteAliasParams{UserID: userID, LineKey: string(key)}))
}
