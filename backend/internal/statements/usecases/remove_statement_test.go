package usecases

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

func TestRemoveStatement_DeletesThePublicID(t *testing.T) {
	statements := &fakeStatements{}

	if err := NewRemoveStatement(statements).Execute(context.Background(), "stmt-1"); err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}

	if !reflect.DeepEqual(statements.deleted, []string{"stmt-1"}) {
		t.Errorf("deleted %v, want [stmt-1]", statements.deleted)
	}
}

func TestRemoveStatement_UnknownPublicIDIsNotFound(t *testing.T) {
	statements := &fakeStatements{deleteErr: shared.ErrNotFound}

	err := NewRemoveStatement(statements).Execute(context.Background(), "stmt-404")

	assertKeyedError(t, err, shared.ErrNotFound, domain.KeyNotFound, nil)
}

func TestRemoveStatement_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")
	statements := &fakeStatements{deleteErr: cause}

	err := NewRemoveStatement(statements).Execute(context.Background(), "stmt-1")

	if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), "delete statement: ") {
		t.Errorf("error = %v, want \"delete statement\" wrapping %v", err, cause)
	}
}
