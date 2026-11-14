package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

var _ ports.StatementRepository = Statements{}

var (
	billingDate = time.Date(2026, 10, 25, 3, 0, 0, 0, time.UTC) // 00:00 in Santiago
	importedAt  = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
)

// importOf builds a statement of the user with one charge per amount.
func importOf(userID int64, publicID string, amounts ...shared.Money) domain.ImportedStatement {
	statement := &domain.Statement{PublicID: publicID, UserID: userID, Bank: "bci", Filename: "cartola.pdf", Kind: domain.FileKindPDF, BillingDate: billingDate, ImportedAt: importedAt, BadLines: []int{4, 9}}
	charges := make([]shared.Charge, 0, len(amounts))
	for i, amount := range amounts {
		charges = append(charges, shared.Charge{
			UserID: userID, Bank: "bci", Source: shared.SourceStatement, ExternalID: publicID + ":" + string(rune('1'+i)),
			OccurredAt: billingDate.AddDate(0, 0, -i), Amount: amount, Currency: shared.CurrencyCLP, Merchant: "NETFLIX", Card: "credito", CardLast4: "1234",
		})
	}
	return domain.ImportedStatement{Statement: statement, Charges: charges}
}

// owners opens a database with two users and returns their ids.
func owners(t *testing.T) (*sql.DB, int64, int64) {
	t.Helper()
	db := sqlitetest.OpenMigrated(t)
	users := sqlitetest.Users{DB: db}
	return db, users.Insert(t, "ana@example.com"), users.Insert(t, "beto@example.com")
}

func chargeCount(t *testing.T, db *sql.DB) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM charges`).Scan(&count); err != nil {
		t.Fatalf("count charges: %v", err)
	}
	return count
}

func TestStatements_CreateSavesTheStatementWithItsCharges(t *testing.T) {
	db, anaID, _ := owners(t)
	imported := importOf(anaID, "st-ana", 599000, 1290000, 899000)

	err := NewStatements(db).Create(context.Background(), anaID, imported)

	if err != nil || imported.Statement.ID == 0 || chargeCount(t, db) != 3 {
		t.Errorf("Create err = %v, ID = %d, charges = %d; want the ID assigned and 3 charges", err, imported.Statement.ID, chargeCount(t, db))
	}
}

func TestStatements_CreateIsAllOrNothing(t *testing.T) {
	db, anaID, _ := owners(t)
	imported := importOf(anaID, "st-ana", 599000, -1, 899000) // the database rejects the negative amount

	err := NewStatements(db).Create(context.Background(), anaID, imported)

	statements, _ := NewStatements(db).List(context.Background(), anaID)
	if err == nil || len(statements) != 0 || chargeCount(t, db) != 0 {
		t.Errorf("Create err = %v, statements = %d, charges = %d; want an error and nothing saved", err, len(statements), chargeCount(t, db))
	}
}

func TestStatements_CreateTwiceIsAConflictOnlyForTheSameUser(t *testing.T) {
	db, anaID, betoID := owners(t)
	statements := NewStatements(db)
	if err := statements.Create(context.Background(), anaID, importOf(anaID, "st-ana-1", 599000)); err != nil {
		t.Fatalf("first Create error = %v", err)
	}

	again := statements.Create(context.Background(), anaID, importOf(anaID, "st-ana-2", 599000))
	other := statements.Create(context.Background(), betoID, importOf(betoID, "st-beto", 599000))

	var keyed shared.KeyedError
	if !errors.Is(again, shared.ErrConflict) || !errors.As(again, &keyed) || keyed.Key != domain.KeyAlreadyImported {
		t.Errorf("same user again error = %v, want ErrConflict with %s", again, domain.KeyAlreadyImported)
	}
	if other != nil {
		t.Errorf("another user's statement error = %v, want nil", other)
	}
}

func TestStatements_ListShowsOnlyTheUsersOwn(t *testing.T) {
	db, anaID, betoID := owners(t)
	statements := NewStatements(db)
	anas := importOf(anaID, "st-ana", 599000, 899000)
	for _, imported := range []domain.ImportedStatement{anas, importOf(betoID, "st-beto", 100)} {
		if err := statements.Create(context.Background(), imported.Statement.UserID, imported); err != nil {
			t.Fatalf("Create error = %v", err)
		}
	}

	got, err := statements.List(context.Background(), anaID)

	want := *anas.Statement
	want.Charges = 2
	if err != nil || len(got) != 1 || !reflect.DeepEqual(got[0], want) {
		t.Errorf("List = %+v, %v; want only %+v", got, err, want)
	}
}

func TestStatements_DeleteRemovesItsCharges(t *testing.T) {
	db, anaID, _ := owners(t)
	statements := NewStatements(db)
	if err := statements.Create(context.Background(), anaID, importOf(anaID, "st-ana", 599000, 899000)); err != nil {
		t.Fatalf("Create error = %v", err)
	}
	if before, _ := statements.List(context.Background(), anaID); len(before) != 1 || chargeCount(t, db) != 2 {
		t.Fatalf("before delete: statements = %d, charges = %d; want 1 and 2", len(before), chargeCount(t, db))
	}

	err := statements.Delete(context.Background(), anaID, "st-ana")

	listed, _ := statements.List(context.Background(), anaID)
	if err != nil || len(listed) != 0 || chargeCount(t, db) != 0 {
		t.Errorf("Delete err = %v, statements = %d, charges = %d; want both gone", err, len(listed), chargeCount(t, db))
	}
}

func TestStatements_DeleteOfAnotherUsersStatementIsNotFound(t *testing.T) {
	db, anaID, betoID := owners(t)
	statements := NewStatements(db)
	if err := statements.Create(context.Background(), anaID, importOf(anaID, "st-ana", 599000)); err != nil {
		t.Fatalf("Create error = %v", err)
	}

	for _, publicID := range []string{"st-ana", "never-imported"} {
		if err := statements.Delete(context.Background(), betoID, publicID); !errors.Is(err, shared.ErrNotFound) {
			t.Errorf("Delete(%s) by another user error = %v, want ErrNotFound", publicID, err)
		}
	}
	if listed, _ := statements.List(context.Background(), anaID); len(listed) != 1 || chargeCount(t, db) != 1 {
		t.Errorf("Ana's statement or charges were touched")
	}
}
