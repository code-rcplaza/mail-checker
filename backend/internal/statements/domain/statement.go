package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Statement is an imported card statement.
type Statement struct {
	ID          int64  // internal
	PublicID    string // UUID v7, the one exposed
	Bank        shared.Bank
	Filename    string
	Kind        FileKind
	BillingDate time.Time
	ImportedAt  time.Time
	Charges     int
	BadLines    []int // line numbers; the back never sends text
}

type FileKind string

const (
	FileKindPDF  FileKind = "pdf"
	FileKindXLSX FileKind = "xlsx"
	FileKindCSV  FileKind = "csv"
	FileKindText FileKind = "text"
)

func ParseFileKind(raw string) (FileKind, error) {
	kind := FileKind(raw)
	switch kind {
	case FileKindPDF, FileKindXLSX, FileKindCSV, FileKindText:
		return kind, nil
	}
	return "", shared.Invalid(KeyFileKindUnknown, nil)
}
