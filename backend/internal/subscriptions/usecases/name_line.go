package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
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

func (u NameLine) Execute(ctx context.Context, key domain.GroupKey, rawAlias string) error {
	alias, err := domain.ParseLineAlias(rawAlias)
	if err != nil {
		return err
	}
	if err := u.ensureLine(ctx, key); err != nil {
		return err
	}
	return u.storeAlias(ctx, key, alias)
}

// ensureLine looks over the widest period, so any line the user sees in the
// period they picked can be named.
func (u NameLine) ensureLine(ctx context.Context, key domain.GroupKey) error {
	groups, err := readGroupsSince(ctx, u.charges, u.catalog, domain.PeriodStart(u.now(), domain.MaxPeriodMonths))
	if err != nil {
		return err
	}
	if !domain.IsLine(groups, key) {
		return shared.NotFound(domain.KeyLineNotFound)
	}
	return nil
}

func (u NameLine) storeAlias(ctx context.Context, key domain.GroupKey, alias string) error {
	if alias == "" {
		return u.deleteAlias(ctx, key)
	}
	if err := u.aliases.Save(ctx, key, alias); err != nil {
		return fmt.Errorf("save alias for %s: %w", key, err)
	}
	return nil
}

func (u NameLine) deleteAlias(ctx context.Context, key domain.GroupKey) error {
	if err := u.aliases.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete alias for %s: %w", key, err)
	}
	return nil
}
