package domain

import (
	"regexp"
	"testing"
)

var featureReasonFormat = regexp.MustCompile(`^[a-z]+(_[a-z]+)*\.[a-z]+(_[a-z]+)*$`)

func TestCommonKeys_FollowFeatureReasonFormat(t *testing.T) {
	cases := []struct {
		key  MessageKey
		want MessageKey
	}{
		{KeyInvalidJSON, "common.invalid_json"},
		{KeyNotFound, "common.not_found"},
		{KeyConflict, "common.conflict"},
		{KeyRouteNotFound, "common.route_not_found"},
		{KeyValidationFailed, "common.validation_failed"},
		{KeyInternal, "common.internal"},
		{KeyCurrencyUnknown, "common.currency_unknown"},
		{KeyAmountNegative, "common.amount_negative"},
		{KeyDateMissing, "common.date_missing"},
		{KeyMerchantMissing, "common.merchant_missing"},
		{KeyExternalIDMissing, "common.external_id_missing"},
		{KeyUserMissing, "common.user_missing"},
		{KeyUnauthorized, "common.unauthorized"},
		{KeyTooManyRequests, "common.too_many_requests"},
		{KeyForbidden, "common.forbidden"},
		{KeyInvalidRequest, "common.invalid_request"},
		{KeyUnavailable, "common.unavailable"},
	}
	for _, tc := range cases {
		t.Run(string(tc.want), func(t *testing.T) {
			if tc.key != tc.want {
				t.Errorf("key = %q, want %q", tc.key, tc.want)
			}
			if !featureReasonFormat.MatchString(string(tc.key)) {
				t.Errorf("key %q does not match <feature>.<reason> in snake_case", tc.key)
			}
		})
	}
}
