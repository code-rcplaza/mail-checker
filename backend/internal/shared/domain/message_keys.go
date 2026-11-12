package domain

// MessageKey is an i18n key in <feature>.<reason> format (English, snake_case).
// The front translates it; the back never sends user-facing text.
type MessageKey string

const (
	KeyInvalidJSON           MessageKey = "common.invalid_json"
	KeyNotFound              MessageKey = "common.not_found"
	KeyConflict              MessageKey = "common.conflict"
	KeyRouteNotFound         MessageKey = "common.route_not_found"
	KeyValidationFailed      MessageKey = "common.validation_failed"
	KeyInternal              MessageKey = "common.internal"
	KeyCurrencyUnknown       MessageKey = "common.currency_unknown"
	KeyAmountNegative        MessageKey = "common.amount_negative"
	KeyDateMissing           MessageKey = "common.date_missing"
	KeyMerchantMissing       MessageKey = "common.merchant_missing"
	KeyExternalIDMissing     MessageKey = "common.external_id_missing"
	KeyCardLastDigitsInvalid MessageKey = "common.card_last_digits_invalid"
)
