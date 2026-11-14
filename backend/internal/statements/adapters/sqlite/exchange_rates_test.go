package sqlite

import (
	"context"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
)

var _ ports.ExchangeRateRecorder = ExchangeRates{}

func TestExchangeRates_RecordOncePerUserAndDate(t *testing.T) {
	db, anaID, betoID := owners(t)
	rates := NewExchangeRates(db)
	ctx := context.Background()
	date := time.Date(2026, 10, 25, 3, 0, 0, 0, time.UTC)
	for _, record := range []struct {
		userID int64
		rate   shared.ExchangeRate
	}{
		{anaID, shared.ExchangeRate{PesosPerDollar: 95012, Date: date}},
		{anaID, shared.ExchangeRate{PesosPerDollar: 99999, Date: date}},
		{betoID, shared.ExchangeRate{PesosPerDollar: 96000, Date: date}},
	} {
		if err := rates.Record(ctx, record.userID, record.rate); err != nil {
			t.Fatalf("Record error = %v, want nil (idempotent)", err)
		}
	}

	stored := map[int64]int64{}
	rows, err := db.Query(`SELECT user_id, pesos_per_dollar FROM exchange_rates`)
	if err != nil {
		t.Fatalf("read rates: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var userID, pesos int64
		if err := rows.Scan(&userID, &pesos); err != nil {
			t.Fatalf("scan rate: %v", err)
		}
		stored[userID] = pesos
	}
	if len(stored) != 2 || stored[anaID] != 95012 || stored[betoID] != 96000 {
		t.Errorf("stored rates = %v, want Ana's first rate and Beto's, one each", stored)
	}
}
