// Package sqlite stores the auth feature in the app database.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite/queries"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	sharedsqlite "github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

type Users struct {
	queries *queries.Queries
	now     func() time.Time
}

func NewUsers(db *sql.DB, now func() time.Time) Users {
	return Users{queries: queries.New(db), now: now}
}

// Create answers ErrConflict for a taken email, whatever its casing.
func (u Users) Create(ctx context.Context, user domain.User) (domain.User, error) {
	id, err := u.queries.CreateUser(ctx, queries.CreateUserParams{
		PublicID:        user.PublicID,
		Name:            user.Name,
		Email:           user.Email,
		AcceptedTermsAt: user.AcceptedTermsAt.Unix(),
		CreatedAt:       u.now().Unix(),
	})
	if sharedsqlite.IsUniqueViolation(err) {
		return domain.User{}, shared.Conflict(shared.KeyConflict)
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	user.ID = id
	user.AcceptedTermsAt = sharedsqlite.TimeFromUnix(user.AcceptedTermsAt.Unix())
	return user, nil
}

func (u Users) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := u.queries.FindUserByEmail(ctx, email)
	return userFromRow(row, err)
}

func (u Users) FindByID(ctx context.Context, id int64) (domain.User, error) {
	row, err := u.queries.FindUserByID(ctx, id)
	return userFromRow(row, err)
}

func userFromRow(row queries.User, err error) (domain.User, error) {
	if err != nil {
		return domain.User{}, lookupError("find user", err)
	}
	return domain.User{
		ID:              row.ID,
		PublicID:        row.PublicID,
		Name:            row.Name,
		Email:           row.Email,
		AcceptedTermsAt: sharedsqlite.TimeFromUnix(row.AcceptedTermsAt),
	}, nil
}

// lookupError turns "no row" into ErrNotFound, the answer every port promises.
func lookupError(step string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", step, shared.ErrNotFound)
	}
	return fmt.Errorf("%s: %w", step, err)
}
