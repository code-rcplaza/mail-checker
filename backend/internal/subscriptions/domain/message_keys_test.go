package domain

import (
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestSubscriptionsKeysExist(t *testing.T) {
	cases := []struct {
		key  shared.MessageKey
		want shared.MessageKey
	}{
		{KeyCancelURLInvalid, "subscriptions.cancel_url_invalid"},
		{KeyRiseNotFound, "subscriptions.rise_not_found"},
		{KeyManualNameMissing, "subscriptions.manual_name_missing"},
		{KeyManualNameTooLong, "subscriptions.manual_name_too_long"},
		{KeyManualAmountNotPositive, "subscriptions.manual_amount_not_positive"},
		{KeyManualDayOutOfRange, "subscriptions.manual_day_out_of_range"},
		{KeyManualNotFound, "subscriptions.manual_not_found"},
		{KeyOverrideAmountNotPositive, "subscriptions.override_amount_not_positive"},
		{KeyOverrideDayOutOfRange, "subscriptions.override_day_out_of_range"},
		{KeyLineKeyRequired, "subscriptions.line_key_required"},
	}
	for _, tc := range cases {
		if tc.key != tc.want {
			t.Errorf("key = %q, want %q", tc.key, tc.want)
		}
	}
}
