// Package accounts bridges emailsync to the accounts feature; it is the only
// place where one feature imports the other.
package accounts

import (
	"context"
	"fmt"

	accountsdomain "github.com/code-rcplaza/mail-checker/backend/internal/accounts/domain"
	accountsusecases "github.com/code-rcplaza/mail-checker/backend/internal/accounts/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
)

// AccountDirectory implements emailsync's port with the accounts use cases.
type AccountDirectory struct {
	credentials accountsusecases.Credentials
	recordSync  accountsusecases.RecordSync
}

func NewAccountDirectory(credentials accountsusecases.Credentials, recordSync accountsusecases.RecordSync) AccountDirectory {
	return AccountDirectory{credentials: credentials, recordSync: recordSync}
}

func (d AccountDirectory) Credentials(ctx context.Context) ([]domain.Credential, error) {
	credentials, err := d.credentials.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("read account credentials: %w", err)
	}
	return toEmailsyncCredentials(credentials), nil
}

func (d AccountDirectory) RecordSync(ctx context.Context, record domain.SyncRecord) error {
	return d.recordSync.Execute(ctx, accountsdomain.SyncRecord{AccountID: record.AccountID, At: record.At, ErrorKey: record.ErrorKey})
}

func toEmailsyncCredentials(credentials []accountsdomain.Credential) []domain.Credential {
	converted := make([]domain.Credential, 0, len(credentials))
	for _, credential := range credentials {
		converted = append(converted, domain.Credential{
			AccountID:       credential.AccountID,
			AccountPublicID: credential.AccountPublicID,
			Email:           credential.Email,
			Password:        credential.Password,
		})
	}
	return converted
}
