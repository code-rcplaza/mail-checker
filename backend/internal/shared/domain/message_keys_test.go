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
