package usecases

import (
	"context"
	"fmt"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

// readChargesSince is shared by both views so they fail the same way.
func readChargesSince(ctx context.Context, reader ports.ChargeReader, since time.Time) ([]shared.Charge, error) {
	charges, err := reader.ListSince(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("read charges since %s: %w", since.Format(time.DateOnly), err)
	}
	return charges, nil
}
