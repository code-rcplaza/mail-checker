package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

// DeleteUser removes the account and everything of the user; their sessions
// go with it, so no cookie opens anything afterwards.
type DeleteUser struct {
	users ports.UserRepository
}

func NewDeleteUser(users ports.UserRepository) DeleteUser { return DeleteUser{users: users} }

func (u DeleteUser) Execute(ctx context.Context, userID int64) error {
	if err := u.users.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	return nil
}
