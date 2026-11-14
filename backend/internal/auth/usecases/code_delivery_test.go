package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestCodeDelivery_IssueCodeStoresOnlyTheHash(t *testing.T) {
	fixture := newDeliveryFixture()

	code, err := fixture.delivery().issueCode(context.Background(), ana.ID)

	stored := fixture.codes.codes[ana.ID]
	if err != nil || len(code) != domain.CodeDigits || !stored.Matches(code) || stored.CodeHash == code {
		t.Errorf("code = %q, stored = %+v, err = %v; want the clear code returned and only its hash stored", code, stored, err)
	}
}

func TestCodeDelivery_ANewCodeReplacesThePrevious(t *testing.T) {
	fixture := newDeliveryFixture()
	fixture.codes.codes[ana.ID] = domain.LoginCode{UserID: ana.ID, CodeHash: domain.HashSecret("111111"), ExpiresAt: testNow.Add(time.Minute), Attempts: 2}

	code, err := fixture.delivery().issueCode(context.Background(), ana.ID)

	stored := fixture.codes.codes[ana.ID]
	if err != nil || stored.Matches("111111") || !stored.Matches(code) || stored.Attempts != 0 {
		t.Errorf("stored = %+v, err = %v; want only the new code %q, without attempts", stored, err, code)
	}
}

func TestCodeDelivery_IssueCodeWrapsSaveError(t *testing.T) {
	fixture := newDeliveryFixture()
	fixture.codes.saveErr = errDatabase

	code, err := fixture.delivery().issueCode(context.Background(), ana.ID)

	assertWrapped(t, err, errDatabase)
	if code != "" {
		t.Errorf("code = %q, want none", code)
	}
}

func TestCodeDelivery_AdmitRequestLimitsPerEmail(t *testing.T) {
	recent := []time.Time{testNow.Add(-50 * time.Minute), testNow.Add(-20 * time.Minute), testNow.Add(-5 * time.Minute)}
	cases := []struct {
		name      string
		email     string
		wantErr   bool
		wantSends int
	}{
		{"three recent sends block another", "ana@example.com", true, 3},
		{"another email is not limited", "beto@example.com", false, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newDeliveryFixture()
			fixture.sends["ana@example.com"] = recent

			err := fixture.delivery().admitRequest(context.Background(), tc.email)

			if tc.wantErr {
				assertKeyedError(t, err, shared.ErrTooMany, domain.KeyTooManySends)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("admitRequest error = %v, want nil", err)
			}
			if got := len(fixture.sends[tc.email]); got != tc.wantSends {
				t.Errorf("sends of %s = %d, want %d", tc.email, got, tc.wantSends)
			}
		})
	}
}

func TestCodeDelivery_AnAdmittedRequestIsRecordedAtNow(t *testing.T) {
	fixture := newDeliveryFixture()

	err := fixture.delivery().admitRequest(context.Background(), "ana@example.com")

	if got := fixture.sends["ana@example.com"]; err != nil || len(got) != 1 || !got[0].Equal(testNow) {
		t.Errorf("sends = %v, err = %v; want one send at %v", got, err, testNow)
	}
}
