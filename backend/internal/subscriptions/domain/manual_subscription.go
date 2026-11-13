package domain

import (
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

const (
	KeyManualNameMissing       shared.MessageKey = "subscriptions.manual_name_missing"
	KeyManualNameTooLong       shared.MessageKey = "subscriptions.manual_name_too_long"
	KeyManualAmountNotPositive shared.MessageKey = "subscriptions.manual_amount_not_positive"
	KeyManualDayOutOfRange     shared.MessageKey = "subscriptions.manual_day_out_of_range"
	KeyManualNotFound          shared.MessageKey = "subscriptions.manual_not_found"
)

// Billing days a month can have; a day the month lacks lands on its last day.
const (
	firstBillingDay = 1
	lastBillingDay  = 31
)

// manualKeyPrefix marks the line of a manual subscription, which has no
// charges and so no parent from the catalog.
const manualKeyPrefix = "manual:"

// ManualSubscription is one the user added by hand because it is not in
// their statements. It is monthly and adds to the spend like any other.
type ManualSubscription struct {
	PublicID   string
	Name       string
	Amount     shared.Money
	Currency   shared.Currency
	BillingDay int
}

// NewManualSubscription validates what the user typed. The name is as long
// as a line alias, since it shows in the same place.
func NewManualSubscription(publicID, name string, amount shared.Money, currency string, billingDay int) (ManualSubscription, error) {
	trimmed := strings.TrimSpace(name)
	if err := validateManualName(trimmed); err != nil {
		return ManualSubscription{}, err
	}
	if amount <= 0 {
		return ManualSubscription{}, shared.Invalid(KeyManualAmountNotPositive, nil)
	}
	if billingDay < firstBillingDay || billingDay > lastBillingDay {
		return ManualSubscription{}, shared.Invalid(KeyManualDayOutOfRange, nil)
	}
	parsedCurrency, err := shared.ParseCurrency(currency)
	if err != nil {
		return ManualSubscription{}, err
	}
	return ManualSubscription{PublicID: publicID, Name: trimmed, Amount: amount, Currency: parsedCurrency, BillingDay: billingDay}, nil
}

func validateManualName(name string) error {
	if name == "" {
		return shared.Invalid(KeyManualNameMissing, nil)
	}
	if utf8.RuneCountInString(name) > MaxLineAliasRunes {
		return shared.Invalid(KeyManualNameTooLong, map[string]any{"max": MaxLineAliasRunes})
	}
	return nil
}

// NextManualCharge is the next time the billing day comes, today included,
// in the time zone now carries.
func NextManualCharge(billingDay int, now time.Time) time.Time {
	thisMonth := dayOfMonth(now.Year(), now.Month(), billingDay, now.Location())
	if thisMonth.Day() >= now.Day() {
		return thisMonth
	}
	return dayOfMonth(now.Year(), now.Month()+1, billingDay, now.Location())
}

// AddManualSubscriptions puts the manual ones in Subscriptions and counts
// them in the monthly spend and the charges still to come this month; the
// given summary is left untouched.
func AddManualSubscriptions(summary Summary, manuals []ManualSubscription, now time.Time) Summary {
	subscriptions := slices.Clone(summary.Subscriptions)
	for _, manual := range manuals {
		subscriptions = append(subscriptions, manualLine(manual, now))
	}
	summary.Subscriptions = subscriptions
	summary.MonthlySpend = MonthlySpend(subscriptions)
	summary.RemainingThisMonth = CountRemainingThisMonth(subscriptions, now)
	return summary
}

func manualLine(manual ManualSubscription, now time.Time) Subscription {
	key := GroupKey(manualKeyPrefix + manual.PublicID)
	return Subscription{
		Key:        key,
		Parent:     key,
		Alias:      manual.Name,
		Status:     StatusActive,
		Cadence:    CadenceMonthly,
		Amount:     manual.Amount,
		Currency:   manual.Currency,
		NextCharge: NextManualCharge(manual.BillingDay, now),
	}
}
