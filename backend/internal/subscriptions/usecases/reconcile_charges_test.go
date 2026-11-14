package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

func reconcile(reader *fakeChargeReader, months int) (domain.Reconciliation, error) {
	useCase := NewReconcileCharges(linesOf(reader))
	return useCase.Execute(context.Background(), userOne, months)
}

func idsOf(charges []shared.Charge) []string {
	ids := make([]string, 0, len(charges))
	for _, c := range charges {
		ids = append(ids, c.ExternalID)
	}
	return ids
}

func pairIDsOf(pairs []domain.ChargePair) [][2]string {
	ids := make([][2]string, 0, len(pairs))
	for _, pair := range pairs {
		ids = append(ids, [2]string{pair.Statement.ExternalID, pair.Email.ExternalID})
	}
	return ids
}

func TestReconcileCharges(t *testing.T) {
	cases := []struct {
		name           string
		charges        []shared.Charge
		wantNotified   [][2]string
		wantUnnotified []string
		wantEmailOnly  []string
	}{
		{
			name: "statement and email charge that match are notified",
			charges: []shared.Charge{
				chargeDaysAgo("email", "STREAMING CO", shared.SourceEmail, 599000, 10),
				chargeDaysAgo("statement", "STRMNG*CO", shared.SourceStatement, 599000, 8),
			},
			wantNotified:   [][2]string{{"statement", "email"}},
			wantUnnotified: []string{},
			wantEmailOnly:  []string{},
		},
		{
			name: "statement-only charge is unnotified",
			charges: []shared.Charge{
				chargeDaysAgo("statement", "CORNER SHOP", shared.SourceStatement, 800000, 8),
			},
			wantNotified:   [][2]string{},
			wantUnnotified: []string{"statement"},
			wantEmailOnly:  []string{},
		},
		{
			name: "email-only charge is email only",
			charges: []shared.Charge{
				chargeDaysAgo("email", "STREAMING CO", shared.SourceEmail, 599000, 10),
			},
			wantNotified:   [][2]string{},
			wantUnnotified: []string{},
			wantEmailOnly:  []string{"email"},
		},
		{
			name: "two equal statement charges and one email pair one to one",
			charges: []shared.Charge{
				chargeDaysAgo("email", "STREAMING CO", shared.SourceEmail, 599000, 10),
				chargeDaysAgo("first", "STRMNG*CO", shared.SourceStatement, 599000, 8),
				chargeDaysAgo("second", "STRMNG*CO", shared.SourceStatement, 599000, 8),
			},
			wantNotified:   [][2]string{{"first", "email"}},
			wantUnnotified: []string{"second"},
			wantEmailOnly:  []string{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reconcile(&fakeChargeReader{charges: tc.charges}, 6)

			if err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}
			if pairs := pairIDsOf(got.Notified); !reflect.DeepEqual(pairs, tc.wantNotified) {
				t.Errorf("notified = %v, want %v", pairs, tc.wantNotified)
			}
			if ids := idsOf(got.Unnotified); !reflect.DeepEqual(ids, tc.wantUnnotified) {
				t.Errorf("unnotified = %v, want %v", ids, tc.wantUnnotified)
			}
			if ids := idsOf(got.EmailOnly); !reflect.DeepEqual(ids, tc.wantEmailOnly) {
				t.Errorf("email only = %v, want %v", ids, tc.wantEmailOnly)
			}
		})
	}
}

// The same errors as GetSummary: compared message by message.
func TestReconcileCharges_FailsLikeGetSummary(t *testing.T) {
	errDatabase := errors.New("database is locked")
	cases := []struct {
		name   string
		reader func() *fakeChargeReader
		months int
	}{
		{"months out of range", func() *fakeChargeReader { return &fakeChargeReader{} }, 25},
		{"charge reader fails", func() *fakeChargeReader { return &fakeChargeReader{err: errDatabase} }, 6},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, summaryErr := getSummary(tc.reader(), decisionsByUser{}, tc.months)
			_, reconcileErr := reconcile(tc.reader(), tc.months)

			if reconcileErr == nil || summaryErr == nil {
				t.Fatalf("errors = (%v, %v), want both non-nil", reconcileErr, summaryErr)
			}
			if reconcileErr.Error() != summaryErr.Error() {
				t.Errorf("ReconcileCharges error = %q, want the same as GetSummary %q", reconcileErr, summaryErr)
			}
		})
	}
}
