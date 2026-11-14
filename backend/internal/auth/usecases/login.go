package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type Login struct {
	users    ports.UserRepository
	delivery CodeDelivery
}

func NewLogin(users ports.UserRepository, delivery CodeDelivery) Login {
	return Login{users: users, delivery: delivery}
}

// Execute answers the same for a registered email and an unknown one; only the
// registered one gets a code.
func (u Login) Execute(ctx context.Context, rawEmail string) error {
	email, err := domain.ParseEmail(rawEmail)
	if err != nil {
		return err
	}
	if err := u.delivery.admitRequest(ctx, email); err != nil {
		return err
	}
	user, isRegistered, err := u.findUser(ctx, email)
	if err != nil || !isRegistered {
		return err
	}
	return u.delivery.sendCode(ctx, user)
}

func (u Login) findUser(ctx context.Context, email string) (domain.User, bool, error) {
	user, err := u.users.FindByEmail(ctx, email)
	if errors.Is(err, shared.ErrNotFound) {
		return domain.User{}, false, nil
	}
	if err != nil {
		return domain.User{}, false, fmt.Errorf("find user: %w", err)
	}
	return user, true, nil
}
