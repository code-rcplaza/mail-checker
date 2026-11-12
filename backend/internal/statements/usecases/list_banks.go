package usecases

import (
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

// ListBanks tells the bank selector which statements can be imported.
type ListBanks struct {
	parsers []ports.BankStatementParser
}

// NewListBanks copies the bank registry so it stays read-only after main builds it.
func NewListBanks(parsers []ports.BankStatementParser) ListBanks {
	return ListBanks{parsers: slices.Clone(parsers)}
}

func (u ListBanks) Execute() []shared.Bank {
	banks := make([]shared.Bank, 0, len(u.parsers))
	for _, parser := range u.parsers {
		banks = append(banks, parser.Bank())
	}
	slices.Sort(banks)
	return banks
}
