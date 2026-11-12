package usecases

import (
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

func TestListBanks_ReturnsRegisteredBanksSortedAlphabetically(t *testing.T) {
	parsers := []ports.BankStatementParser{
		&fakeStatementParser{bank: "falabella"},
		&fakeStatementParser{bank: "bci"},
		&fakeStatementParser{bank: "santander"},
	}

	got := NewListBanks(parsers).Execute()

	want := []shared.Bank{"bci", "falabella", "santander"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("banks = %v, want %v", got, want)
	}
}
