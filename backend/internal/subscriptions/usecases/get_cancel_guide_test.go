package usecases

import (
	"context"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// catalogWithCancelLink is fakeCatalog with a verified cancel link for Streaming Co.
type catalogWithCancelLink struct{ fakeCatalog }

func (catalogWithCancelLink) Identify(merchant string) (domain.Service, bool) {
	service, isKnown := fakeCatalog{}.Identify(merchant)
	service.CancelURL = "https://streaming.example/cancel"
	return service, isKnown
}

func TestGetCancelGuide(t *testing.T) {
	charges := chargesByUser{userOne: append(streamingOnce(), chargeDaysAgo("o1", "OTHER SHOP", shared.SourceStatement, 100000, 10))}
	cases := []struct {
		name    string
		catalog ChargeLines
		userID  int64
		key     domain.GroupKey
		want    string
		wantErr error
	}{
		{"a service with a link gives it", ChargeLines{Charges: charges, Catalog: catalogWithCancelLink{}, Now: fixedNow}, userOne, streamingLine, "https://streaming.example/cancel", nil},
		{"a service without a link gives none", linesOf(charges), userOne, streamingLine, "", nil},
		{"an unknown merchant gives none", linesOf(charges), userOne, "merchant:OTHER SHOP:CLP||d02", "", nil},
		{"a key that is not a line is not found", linesOf(charges), userOne, "service:Nothing||d01", "", shared.ErrNotFound},
		{"the line of another user is not found", linesOf(charges), userTwo, streamingLine, "", shared.ErrNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewGetCancelGuide(tc.catalog).Execute(context.Background(), tc.userID, tc.key)

			if got != tc.want {
				t.Errorf("url = %q, want %q", got, tc.want)
			}
			if tc.wantErr != nil {
				assertKeyedError(t, err, tc.wantErr, domain.KeyLineNotFound)
			}
			if tc.wantErr == nil && err != nil {
				t.Errorf("Execute error = %v, want nil", err)
			}
		})
	}
}
