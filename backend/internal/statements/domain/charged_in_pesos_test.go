package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestIsChargedInPesos(t *testing.T) {
	cases := []struct {
		name    string
		origin  shared.Money
		dollars shared.Money
		want    bool
	}{
		{"origin a thousand times the dollars is pesos", 1000000, 1000, true},
		{"origin equal to the dollars is dollars", 2000, 2000, false},
		{"euros close to the dollars are not pesos", 1500, 1620, false},
		{"exactly the minimum ratio counts as pesos", 100000, 1000, true},
		{"one hundredth under the minimum ratio is not pesos", 99999, 1000, false},
		{"zero dollars with an origin is pesos", 599000, 0, true},
		{"origin and dollars in zero are not pesos", 0, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsChargedInPesos(tc.origin, tc.dollars)

			if got != tc.want {
				t.Errorf("IsChargedInPesos(%d, %d) = %v, want %v", tc.origin, tc.dollars, got, tc.want)
			}
		})
	}
}
