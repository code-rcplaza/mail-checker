// Package export is the bridge that gathers the user's data from every
// feature for their copy. It is the only place in auth that imports the use
// cases of other features.
package export

import (
	"context"
	"fmt"

	authdomain "github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	authusecases "github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	statementsdomain "github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	statementsusecases "github.com/code-rcplaza/mail-checker/backend/internal/statements/usecases"
	subscriptionsdomain "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	subscriptionsusecases "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

// UserCopy is everything of one user: what they stored and the subscriptions
// of their summary, the table the CSV carries.
type UserCopy struct {
	Profile       authdomain.User
	Statements    []statementsdomain.Statement
	Saved         subscriptionsusecases.SavedChoices
	Subscriptions []subscriptionsdomain.Subscription
}

type CopySources struct {
	Me         authusecases.GetMe
	Statements statementsusecases.ListStatements
	Saved      subscriptionsusecases.ListSavedChoices
	Summary    subscriptionsusecases.GetSummary
}

type UserCopies struct {
	sources CopySources
}

func NewUserCopies(sources CopySources) UserCopies { return UserCopies{sources: sources} }

// Gather reads the user's copy; each source only ever reads that user.
func (c UserCopies) Gather(ctx context.Context, userID int64) (UserCopy, error) {
	profile, err := c.sources.Me.Execute(ctx, userID)
	if err != nil {
		return UserCopy{}, fmt.Errorf("copy profile: %w", err)
	}
	statements, err := c.sources.Statements.Execute(ctx, userID)
	if err != nil {
		return UserCopy{}, fmt.Errorf("copy statements: %w", err)
	}
	saved, err := c.sources.Saved.Execute(ctx, userID)
	if err != nil {
		return UserCopy{}, fmt.Errorf("copy saved choices: %w", err)
	}
	summary, err := c.sources.Summary.Execute(ctx, userID, subscriptionsdomain.DefaultPeriodMonths)
	if err != nil {
		return UserCopy{}, fmt.Errorf("copy summary: %w", err)
	}
	return UserCopy{Profile: profile, Statements: statements, Saved: saved, Subscriptions: summary.Subscriptions}, nil
}
