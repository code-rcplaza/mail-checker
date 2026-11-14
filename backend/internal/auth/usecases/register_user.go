package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type RegisterUserDependencies struct {
	Users       ports.UserRepository
	NewPublicID func() string
	Now         func() time.Time
}

type RegisterUser struct {
	dependencies RegisterUserDependencies
}

func NewRegisterUser(dependencies RegisterUserDependencies) RegisterUser {
	return RegisterUser{dependencies: dependencies}
}

// Execute answers the same for a new email and a registered one: the
// registered one keeps its user and nothing else is saved.
func (u RegisterUser) Execute(ctx context.Context, form domain.RegistrationForm) (domain.User, error) {
	registration, err := domain.ParseRegistration(form)
	if err != nil {
		return domain.User{}, err
	}
	return u.ensureUser(ctx, registration)
}

func (u RegisterUser) ensureUser(ctx context.Context, registration domain.Registration) (domain.User, error) {
	existing, err := u.dependencies.Users.FindByEmail(ctx, registration.Email)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, shared.ErrNotFound) {
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}
	return u.createUser(ctx, registration)
}

func (u RegisterUser) createUser(ctx context.Context, registration domain.Registration) (domain.User, error) {
	user := domain.User{
		PublicID:        u.dependencies.NewPublicID(),
		Name:            registration.Name,
		Email:           registration.Email,
		AcceptedTermsAt: u.dependencies.Now(),
	}
	created, err := u.dependencies.Users.Create(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return created, nil
}
