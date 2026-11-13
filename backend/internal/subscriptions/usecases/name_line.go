package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// NameLine stores the user's name for a line (Personal, Padres…) so parallel
// subscriptions of one service can be told apart. A blank name clears it.
type NameLine struct {
	charges ports.ChargeReader
	catalog ports.ServiceCatalog
	aliases ports.LineAliases
	now     func() time.Time
}

func NewNameLine(
	charges ports.ChargeReader,
	catalog ports.ServiceCatalog,
	aliases ports.LineAliases,
	now func() time.Time,
) NameLine {
	return NameLine{charges: charges, catalog: catalog, aliases: aliases, now: now}
}

func (u NameLine) Execute(ctx context.Context, userID int64, key domain.GroupKey, rawAlias string) error {
	alias, err := domain.ParseLineAlias(rawAlias)
	if err != nil {
		return err
	}
	if err := u.ensureLine(ctx, userID, key); err != nil {
		return err
	}
	return u.storeAlias(ctx, userID, key, alias)
}

func (u NameLine) ensureLine(ctx context.Context, userID int64, key domain.GroupKey) error {
	return ensureLineExists(ctx, userID, u.charges, u.catalog, u.now(), key)
}

func (u NameLine) storeAlias(ctx context.Context, userID int64, key domain.GroupKey, alias string) error {
	if alias == "" {
		return u.deleteAlias(ctx, userID, key)
	}
	if err := u.aliases.Save(ctx, userID, key, alias); err != nil {
		return fmt.Errorf("save alias for %s: %w", key, err)
	}
	return nil
}

func (u NameLine) deleteAlias(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.aliases.Delete(ctx, userID, key); err != nil {
		return fmt.Errorf("delete alias for %s: %w", key, err)
	}
	return nil
}
