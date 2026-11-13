package usecases

import (
	"context"
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

const (
	userOne int64 = 1
	userTwo int64 = 2
)

func importAs(t *testing.T, fixture importFixture, userID int64) domain.Statement {
	t.Helper()
	result, err := fixture.useCase.Execute(context.Background(), userID, bciInput())
	if err != nil {
		t.Fatalf("Execute for user %d error = %v, want nil", userID, err)
	}
	return result.Statement
}

func TestFakeStatements_ListsOnlyTheStatementsOfTheUser(t *testing.T) {
	fixture := newImportFixture()
	mine := importAs(t, fixture, userOne)
	importAs(t, fixture, userTwo)

	got, err := fixture.statements.List(context.Background(), userOne)

	if err != nil || len(got) != 1 || got[0].PublicID != mine.PublicID || got[0].UserID != userOne {
		t.Errorf("List = (%+v, %v), want only the statement of user %d", got, err, userOne)
	}
}

func TestFakeStatements_DeleteOfAnotherUsersStatementIsNotFound(t *testing.T) {
	fixture := newImportFixture()
	theirs := importAs(t, fixture, userTwo)

	err := fixture.statements.Delete(context.Background(), userOne, theirs.PublicID)

	if !errors.Is(err, shared.ErrNotFound) {
		t.Errorf("Delete error = %v, want ErrNotFound", err)
	}
}

func TestImportStatement_StatementAndChargesBelongToTheUser(t *testing.T) {
	fixture := newImportFixture()

	got := importAs(t, fixture, userTwo)

	if got.UserID != userTwo {
		t.Errorf("statement user = %d, want %d", got.UserID, userTwo)
	}
	for _, c := range fixture.statements.created[0].charges {
		if c.UserID != userTwo {
			t.Errorf("charge %s user = %d, want %d", c.ExternalID, c.UserID, userTwo)
		}
	}
}

func TestImportStatement_SameFileTwice(t *testing.T) {
	cases := []struct {
		name       string
		secondUser int64
		wantErr    error
	}{
		{"by two users there is no conflict", userTwo, nil},
		{"by the same user it is already imported", userOne, shared.ErrConflict},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newImportFixture()
			importAs(t, fixture, userOne)

			_, err := fixture.useCase.Execute(context.Background(), tc.secondUser, bciInput())

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("second import error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				assertKeyedError(t, err, shared.ErrConflict, domain.KeyAlreadyImported, nil)
			}
		})
	}
}
