package sqlite

import (
	"context"
	"database/sql"
	"time"

	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type HiddenLines struct {
	queries *queries.Queries
}

func NewHiddenLines(db *sql.DB) HiddenLines { return HiddenLines{queries: queries.New(db)} }

func (h HiddenLines) List(ctx context.Context, userID int64) ([]domain.GroupKey, error) {
	rows, err := h.queries.ListHiddenLines(ctx, userID)
	if err != nil {
		return nil, wrapped("list hidden lines", err)
	}
	keys := make([]domain.GroupKey, 0, len(rows))
	for _, key := range rows {
		keys = append(keys, domain.GroupKey(key))
	}
	return keys, nil
}

// Hide keeps one mark per line: hiding again does nothing.
func (h HiddenLines) Hide(ctx context.Context, userID int64, hidden domain.HiddenLine) error {
	return wrapped("hide line", h.queries.HideLine(ctx, queries.HideLineParams{UserID: userID, LineKey: string(hidden.Key), HiddenAt: hidden.HiddenAt.Unix()}))
}

type Cancellations struct {
	queries *queries.Queries
}

func NewCancellations(db *sql.DB) Cancellations { return Cancellations{queries: queries.New(db)} }

func (c Cancellations) List(ctx context.Context, userID int64) (map[domain.GroupKey]time.Time, error) {
	rows, err := c.queries.ListCancellations(ctx, userID)
	if err != nil {
		return nil, wrapped("list cancellations", err)
	}
	marks := make(map[domain.GroupKey]time.Time, len(rows))
	for _, row := range rows {
		marks[domain.GroupKey(row.LineKey)] = sharedsqlite.TimeFromUnix(row.MarkedAt)
	}
	return marks, nil
}

// Mark keeps the date of the first mark: marking again does nothing.
func (c Cancellations) Mark(ctx context.Context, userID int64, mark domain.CancellationMark) error {
	return wrapped("mark for cancellation", c.queries.MarkForCancellation(ctx, queries.MarkForCancellationParams{UserID: userID, LineKey: string(mark.Key), MarkedAt: mark.MarkedAt.Unix()}))
}
