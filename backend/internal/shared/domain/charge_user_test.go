package domain

import (
	"testing"
	"time"
)

func TestNewCharge_UserID(t *testing.T) {
	cases := []struct {
		name    string
		userID  int64
		isValid bool
	}{
		{"owner above zero is valid", 42, true},
		{"zero owner is rejected", 0, false},
		{"negative owner is rejected", -1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := validCharge()
			c.UserID = tc.userID

			got, err := NewCharge(c)

			if !tc.isValid {
				assertInvalidWithKey(t, err, KeyUserMissing)
				return
			}
			if err != nil || got.UserID != tc.userID {
				t.Errorf("NewCharge = (user %d, %v), want user %d and no error", got.UserID, err, tc.userID)
			}
		})
	}
}

func TestToCharge_CarriesTheUserOfTheOrigin(t *testing.T) {
	origin := validOrigin()
	origin.UserID = 42
	parsed := ParsedCharge{Amount: 599000, Currency: CurrencyCLP, Merchant: "STREAMING CO", Date: time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)}

	got, err := parsed.ToCharge(origin)

	if err != nil || got.UserID != 42 {
		t.Errorf("ToCharge = (user %d, %v), want user 42 and no error", got.UserID, err)
	}
}
