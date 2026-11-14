package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// listingUnparsed records the limit it was asked for.
type listingUnparsed struct {
	emails     []domain.UnparsedEmail
	err        error
	askedLimit int
	calls      int
}

func (l *listingUnparsed) Save(context.Context, domain.UnparsedEmail) error { return nil }

func (l *listingUnparsed) Forget(context.Context, string) error { return nil }

func (l *listingUnparsed) List(_ context.Context, limit int) ([]domain.UnparsedEmail, error) {
	l.calls++
	l.askedLimit = limit
	return l.emails, l.err
}

func TestListUnparsed_ReturnsRepositoryEmailsWithTheLimit(t *testing.T) {
	stored := []domain.UnparsedEmail{
		domain.UnparsedFromEmail(domain.BankEmail{Email: purchase("<x1@other>", "newsletter@other.example", "OFFERS"), AccountID: 1, Bank: "bank-a"}),
	}
	repository := &listingUnparsed{emails: stored}

	got, err := NewListUnparsed(repository).Execute(context.Background(), 50)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, stored) {
		t.Errorf("emails = %+v, want %+v", got, stored)
	}
	if repository.askedLimit != 50 {
		t.Errorf("asked limit = %d, want 50", repository.askedLimit)
	}
}

func TestListUnparsed_RejectsLimitOutOfRangeWithoutQuerying(t *testing.T) {
	for _, limit := range []int{0, 101} {
		repository := &listingUnparsed{}

		_, err := NewListUnparsed(repository).Execute(context.Background(), limit)

		assertKeyedError(t, err, shared.ErrInvalid, domain.KeyLimitOutOfRange)
		if repository.calls != 0 {
			t.Errorf("limit = %d: repository queried %d times, want 0", limit, repository.calls)
		}
	}
}

func TestListUnparsed_WrapsRepositoryError(t *testing.T) {
	errDatabase := errors.New("database is locked")

	_, err := NewListUnparsed(&listingUnparsed{err: errDatabase}).Execute(context.Background(), 50)

	if !errors.Is(err, errDatabase) {
		t.Fatalf("errors.Is(err, errDatabase) = false (err = %v), want true", err)
	}
	if err.Error() == errDatabase.Error() {
		t.Errorf("err = %q, want it wrapped with context", err)
	}
}
