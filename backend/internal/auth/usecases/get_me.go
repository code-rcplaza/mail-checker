package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

type GetMe struct {
	users ports.UserRepository
}

func NewGetMe(users ports.UserRepository) GetMe { return GetMe{users: users} }

func (u GetMe) Execute(ctx context.Context, userID int64) (domain.User, error) {
	user, err := u.users.FindByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get me: %w", err)
	}
	return user, nil
}
