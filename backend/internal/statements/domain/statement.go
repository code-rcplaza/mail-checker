package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Statement is an imported card statement.
type Statement struct {
	ID          int64  // internal
	PublicID    string // UUID v7, the one exposed
	UserID      int64  // the owner
	Bank        shared.Bank
	Filename    string
	Kind        FileKind
	BillingDate time.Time
	ImportedAt  time.Time
	Charges     int
	BadLines    []int // line numbers; the back never sends text
}

// ImportedStatement is a statement ready to save with its charges; saving
// assigns the statement's internal ID.
type ImportedStatement struct {
	Statement *Statement
	Charges   []shared.Charge
}

type FileKind string

const (
	FileKindPDF  FileKind = "pdf"
	FileKindXLSX FileKind = "xlsx"
	FileKindCSV  FileKind = "csv"
	FileKindText FileKind = "text"
)
