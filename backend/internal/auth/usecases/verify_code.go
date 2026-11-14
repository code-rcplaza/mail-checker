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

// CodeEntry is what the person types to enter.
type CodeEntry struct {
	Email string
	Code  string
}

type VerifyCodeDependencies struct {
	Users ports.UserRepository
	Codes ports.LoginCodes
	Now   func() time.Time
}

type VerifyCode struct {
	dependencies VerifyCodeDependencies
}

func NewVerifyCode(dependencies VerifyCodeDependencies) VerifyCode {
	return VerifyCode{dependencies: dependencies}
}

// Execute answers auth.code_invalid both for an unknown email and for one
// without a code, so it does not reveal which emails exist.
func (u VerifyCode) Execute(ctx context.Context, entry CodeEntry) (domain.User, error) {
	user, err := u.findUser(ctx, entry.Email)
	if err != nil {
		return domain.User{}, err
	}
	code, err := u.findCode(ctx, user.ID)
	if err != nil {
		return domain.User{}, err
	}
	if code.IsExpiredAt(u.dependencies.Now()) {
		return domain.User{}, shared.Invalid(domain.KeyCodeExpired, nil)
	}
	if !code.Matches(entry.Code) {
		return domain.User{}, u.countFailedAttempt(ctx, user.ID)
	}
	if err := u.consumeCode(ctx, user.ID); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u VerifyCode) findUser(ctx context.Context, rawEmail string) (domain.User, error) {
	email, err := domain.ParseEmail(rawEmail)
	if err != nil {
		return domain.User{}, err
	}
	user, err := u.dependencies.Users.FindByEmail(ctx, email)
	return user, codeInvalidWhenMissing(err, "find user")
}

func (u VerifyCode) findCode(ctx context.Context, userID int64) (domain.LoginCode, error) {
	code, err := u.dependencies.Codes.Find(ctx, userID)
	return code, codeInvalidWhenMissing(err, "find login code")
}

func codeInvalidWhenMissing(err error, step string) error {
	if errors.Is(err, shared.ErrNotFound) {
		return shared.Invalid(domain.KeyCodeInvalid, nil)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", step, err)
	}
	return nil
}

// countFailedAttempt always ends in a rejection: the code is wrong, or that was its last attempt.
func (u VerifyCode) countFailedAttempt(ctx context.Context, userID int64) error {
	attempts, err := u.dependencies.Codes.AddFailedAttempt(ctx, userID)
	if err != nil {
		return fmt.Errorf("count failed attempt: %w", err)
	}
	if !domain.AttemptsExhausted(attempts) {
		return shared.Invalid(domain.KeyCodeInvalid, nil)
	}
	if err := u.consumeCode(ctx, userID); err != nil {
		return err
	}
	return shared.Invalid(domain.KeyCodeAttemptsExceeded, nil)
}

func (u VerifyCode) consumeCode(ctx context.Context, userID int64) error {
	if err := u.dependencies.Codes.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete login code: %w", err)
	}
	return nil
}
