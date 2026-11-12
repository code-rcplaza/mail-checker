package usecases

import (
	"context"
	"fmt"

	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/accounts/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Credentials hands the sync every account with its app password decrypted.
type Credentials struct {
	accounts ports.AccountRepository
	cipher   ports.SecretCipher
}

func NewCredentials(accounts ports.AccountRepository, cipher ports.SecretCipher) Credentials {
	return Credentials{accounts: accounts, cipher: cipher}
}

func (u Credentials) Execute(ctx context.Context) ([]domain.Credential, error) {
	accounts, err := u.readAccounts(ctx)
	if err != nil {
		return nil, err
	}
	return u.decryptAll(accounts)
}

func (u Credentials) readAccounts(ctx context.Context) ([]domain.Account, error) {
	accounts, err := u.accounts.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return accounts, nil
}

// decryptAll fails as a whole: a password sealed with a lost key means the
// user must enter it again, and syncing only some accounts would hide that.
func (u Credentials) decryptAll(accounts []domain.Account) ([]domain.Credential, error) {
	credentials := make([]domain.Credential, 0, len(accounts))
	for _, account := range accounts {
		credential, err := u.decrypt(account)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}

// decrypt drops the cipher's error on purpose: the key the user sees is enough
// and the cause must not risk carrying secret material.
func (u Credentials) decrypt(account domain.Account) (domain.Credential, error) {
	password, err := u.cipher.Decrypt(account.PasswordEnc)
	if err != nil {
		return domain.Credential{}, shared.Conflict(domain.KeyPasswordUnreadable)
	}
	return domain.Credential{
		AccountID:       account.ID,
		AccountPublicID: account.PublicID,
		Email:           account.Email,
		Password:        password,
	}, nil
}
