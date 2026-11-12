package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/ports"
)

type AddAccountInput struct {
	Email       string
	AppPassword string
}

// AddAccount registers a mailbox; the app password only leaves this use case sealed.
type AddAccount struct {
	accounts    ports.AccountRepository
	cipher      ports.SecretCipher
	newPublicID func() string
	now         func() time.Time
}

func NewAddAccount(accounts ports.AccountRepository, cipher ports.SecretCipher, newPublicID func() string, now func() time.Time) AddAccount {
	return AddAccount{accounts: accounts, cipher: cipher, newPublicID: newPublicID, now: now}
}

func (u AddAccount) Execute(ctx context.Context, input AddAccountInput) (domain.Account, error) {
	email, password, err := parseInput(input)
	if err != nil {
		return domain.Account{}, err
	}
	sealed, err := u.encryptPassword(password)
	if err != nil {
		return domain.Account{}, err
	}
	account := u.buildAccount(email, sealed)
	if err := u.saveAccount(ctx, &account); err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

func parseInput(input AddAccountInput) (email, password string, err error) {
	email, err = domain.ParseEmail(input.Email)
	if err != nil {
		return "", "", err
	}
	password, err = domain.ParseAppPassword(input.AppPassword)
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}

func (u AddAccount) encryptPassword(password string) ([]byte, error) {
	sealed, err := u.cipher.Encrypt(password)
	if err != nil {
		return nil, fmt.Errorf("encrypt app password: %w", err)
	}
	return sealed, nil
}

func (u AddAccount) buildAccount(email string, sealed []byte) domain.Account {
	return domain.Account{
		PublicID:    u.newPublicID(),
		Email:       email,
		PasswordEnc: sealed,
		CreatedAt:   u.now(),
	}
}

func (u AddAccount) saveAccount(ctx context.Context, account *domain.Account) error {
	if err := u.accounts.Create(ctx, account); err != nil {
		return fmt.Errorf("save account: %w", err)
	}
	return nil
}
