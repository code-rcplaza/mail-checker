package usecases

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

func TestListStatements_ReturnsTheRepositoryStatements(t *testing.T) {
	saved := []domain.Statement{
		{ID: 1, PublicID: "stmt-1", Bank: "bci", Filename: "cartola-septiembre.pdf", Kind: domain.FileKindPDF, BillingDate: billingDate, ImportedAt: testNow, Charges: 3},
		{ID: 2, PublicID: "stmt-2", Bank: "bci", Filename: "cartola-octubre.csv", Kind: domain.FileKindCSV, BillingDate: billingDate, ImportedAt: testNow, Charges: 1, BadLines: []int{7}},
	}
	statements := &fakeStatements{listed: saved}

	got, err := NewListStatements(statements).Execute(context.Background(), userOne)

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, saved) {
		t.Errorf("listed %+v, want %+v", got, saved)
	}
}

func TestListStatements_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")
	statements := &fakeStatements{listErr: cause}

	_, err := NewListStatements(statements).Execute(context.Background(), userOne)

	if !errors.Is(err, cause) || !strings.HasPrefix(err.Error(), "list statements: ") {
		t.Errorf("error = %v, want \"list statements\" wrapping %v", err, cause)
	}
}
