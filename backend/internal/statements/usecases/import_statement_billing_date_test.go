package usecases

import (
	"context"
	"testing"
	"time"
)

var fileBillingDate = time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)

func TestImportStatement_ChoosesTheBillingDate(t *testing.T) {
	cases := []struct {
		name     string
		fileDate time.Time
		userDate time.Time
		want     time.Time
	}{
		{"file date without a user date", fileBillingDate, time.Time{}, fileBillingDate},
		{"file date wins over the user date", fileBillingDate, billingDate, fileBillingDate},
		{"user date when the file has none", time.Time{}, billingDate, billingDate},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newImportFixture()
			fixture.bci.billingDate = tc.fileDate
			input := bciInput()
			input.BillingDate = tc.userDate

			got, err := fixture.useCase.Execute(context.Background(), userOne, input)

			if err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}
			if !got.Statement.BillingDate.Equal(tc.want) {
				t.Errorf("billing date = %v, want %v", got.Statement.BillingDate, tc.want)
			}
			if saved := fixture.statements.created[0].statement.BillingDate; !saved.Equal(tc.want) {
				t.Errorf("saved billing date = %v, want %v", saved, tc.want)
			}
		})
	}
}
