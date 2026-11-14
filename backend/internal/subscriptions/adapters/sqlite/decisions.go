package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

type Decisions struct {
	queries *queries.Queries
	now     func() time.Time
}

func NewDecisions(db *sql.DB, now func() time.Time) Decisions {
	return Decisions{queries: queries.New(db), now: now}
}

func (d Decisions) List(ctx context.Context, userID int64) (map[domain.GroupKey]domain.Decision, error) {
	rows, err := d.queries.ListDecisions(ctx, userID)
	if err != nil {
		return nil, wrapped("list decisions", err)
	}
	decisions := make(map[domain.GroupKey]domain.Decision, len(rows))
	for _, row := range rows {
		decisions[domain.GroupKey(row.GroupKey)] = domain.Decision{Kind: domain.DecisionKind(row.Decision), Cadence: domain.Cadence(row.Cadence.String)}
	}
	return decisions, nil
}

// Save keeps only the last decision of the line.
func (d Decisions) Save(ctx context.Context, userID int64, decision domain.LineDecision) error {
	return wrapped("save decision", d.queries.SaveDecision(ctx, queries.SaveDecisionParams{
		UserID:    userID,
		GroupKey:  string(decision.Key),
		Decision:  string(decision.Decision.Kind),
		Cadence:   nullWhenEmpty(string(decision.Decision.Cadence)),
		DecidedAt: d.now().Unix(),
	}))
}

func (d Decisions) Delete(ctx context.Context, userID int64, key domain.GroupKey) error {
	return wrapped("delete decision", d.queries.DeleteDecision(ctx, queries.DeleteDecisionParams{UserID: userID, GroupKey: string(key)}))
}
