package sqlite

import (
	"context"
	"reflect"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
)

var _ ports.ExchangeRates = ExchangeRates{}

func TestExchangeRates_ListsOnlyTheUsersRatesByDate(t *testing.T) {
	db, anaID, betoID := twoOwners(t)
	day := func(d int) time.Time { return time.Date(2026, 10, d, 3, 0, 0, 0, time.UTC) }
	insert := `INSERT INTO exchange_rates (user_id, pesos_per_dollar, date) VALUES (?, ?, ?)`
	for _, row := range [][]any{{anaID, 95500, day(25)}, {anaID, 94000, day(5)}, {betoID, 99000, day(30)}, {anaID, 95012, day(15)}} {
		if _, err := db.Exec(insert, row[0], row[1], row[2].(time.Time).Unix()); err != nil {
			t.Fatalf("insert rate: %v", err)
		}
	}
	rates := NewExchangeRates(db)

	got, err := rates.List(context.Background(), anaID)
	none, errNone := rates.List(context.Background(), betoID+1)

	want := []shared.ExchangeRate{{PesosPerDollar: 94000, Date: day(5)}, {PesosPerDollar: 95012, Date: day(15)}, {PesosPerDollar: 95500, Date: day(25)}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("Ana's rates = %v, %v; want %v", got, err, want)
	}
	if errNone != nil || len(none) != 0 {
		t.Errorf("rates of a user without any = %v, %v; want none", none, errNone)
	}
}
