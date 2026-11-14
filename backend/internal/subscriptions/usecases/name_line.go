package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// NameLine stores the user's name for a line (Personal, Padres…) so parallel
// subscriptions of one service can be told apart. A blank name clears it.
type NameLine struct {
	lines   ChargeLines
	aliases ports.LineAliases
}

func NewNameLine(lines ChargeLines, aliases ports.LineAliases) NameLine {
	return NameLine{lines: lines, aliases: aliases}
}

func (u NameLine) Execute(ctx context.Context, userID int64, input domain.LineAlias) error {
	alias, err := domain.ParseLineAlias(input.Alias)
	if err != nil {
		return err
	}
	if err := u.lines.ensureLine(ctx, userID, input.Key); err != nil {
		return err
	}
	return u.storeAlias(ctx, userID, domain.LineAlias{Key: input.Key, Alias: alias})
}

func (u NameLine) storeAlias(ctx context.Context, userID int64, alias domain.LineAlias) error {
	if alias.Alias == "" {
		return u.deleteAlias(ctx, userID, alias.Key)
	}
	if err := u.aliases.Save(ctx, userID, alias); err != nil {
		return fmt.Errorf("save alias for %s: %w", alias.Key, err)
	}
	return nil
}

func (u NameLine) deleteAlias(ctx context.Context, userID int64, key domain.GroupKey) error {
	if err := u.aliases.Delete(ctx, userID, key); err != nil {
		return fmt.Errorf("delete alias for %s: %w", key, err)
	}
	return nil
}
