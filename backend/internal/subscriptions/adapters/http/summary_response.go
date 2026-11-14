package http

import (
	"slices"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// The summary as the front reads it: amounts in hundredths with their
// currency, dates in unix seconds, the line key as an opaque string. Inactive
// lines and other charges are calculated but not shown in the MVP, and no
// internal id leaves the back.
type summaryResponse struct {
	Subscriptions  []subscriptionResponse  `json:"subscriptions"`
	Possible       []possibleResponse      `json:"possible"`
	Total          totalResponse           `json:"total"`
	Charged        int                     `json:"charged"`
	Remaining      int                     `json:"remaining"`
	Rises          []riseResponse          `json:"rises"`
	PendingActions []pendingActionResponse `json:"pendingActions"`
	Upcoming       []upcomingResponse      `json:"upcoming"`
}

type moneyResponse struct {
	Amount   shared.Money    `json:"amount"`
	Currency shared.Currency `json:"currency"`
}

type pesosResponse struct {
	Amount        shared.Money `json:"amount"`
	IsApproximate bool         `json:"isApproximate"`
}

type subscriptionResponse struct {
	LineKey                 domain.GroupKey `json:"lineKey"`
	Name                    string          `json:"name"`
	Alias                   string          `json:"alias"`
	Cadence                 domain.Cadence  `json:"cadence"`
	Amount                  moneyResponse   `json:"amount"`
	InPesos                 *pesosResponse  `json:"inPesos"` // null: dollars without an exchange rate
	NextCharge              *int64          `json:"nextCharge"`
	Bank                    shared.Bank     `json:"bank"`
	Card                    shared.Card     `json:"card"`
	CardLast4               string          `json:"cardLast4"`
	HasRise                 bool            `json:"hasRise"`
	MarkedForCancellationAt *int64          `json:"markedForCancellationAt"`
}

// possibleResponse has no cadence nor next charge: those come once the user decides.
type possibleResponse struct {
	LineKey domain.GroupKey      `json:"lineKey"`
	Name    string               `json:"name"`
	Alias   string               `json:"alias"`
	Amount  moneyResponse        `json:"amount"`
	Charges []seenChargeResponse `json:"charges"`
}

type seenChargeResponse struct {
	Date   int64         `json:"date"`
	Amount moneyResponse `json:"amount"`
}

type totalResponse struct {
	Amount        shared.Money    `json:"amount"`
	Currency      shared.Currency `json:"currency"`
	IsApproximate bool            `json:"isApproximate"`
	IsIncomplete  bool            `json:"isIncomplete"`
}

type riseResponse struct {
	LineKey       domain.GroupKey `json:"lineKey"`
	Before        moneyResponse   `json:"before"`
	After         moneyResponse   `json:"after"`
	Since         int64           `json:"since"`
	MonthlyImpact moneyResponse   `json:"monthlyImpact"`
	YearlyImpact  moneyResponse   `json:"yearlyImpact"`
	IsPlanChange  bool            `json:"isPlanChange"`
}

type pendingActionResponse struct {
	Kind    domain.PendingActionKind `json:"kind"`
	LineKey domain.GroupKey          `json:"lineKey"`
}

type upcomingResponse struct {
	LineKey domain.GroupKey `json:"lineKey"`
	Name    string          `json:"name"`
	Alias   string          `json:"alias"`
	Amount  moneyResponse   `json:"amount"`
	Date    int64           `json:"date"`
}

func summaryResponseOf(summary domain.Summary) summaryResponse {
	return summaryResponse{
		Subscriptions:  subscriptionsOf(summary),
		Possible:       possibleListOf(summary.Possible),
		Total:          totalResponse{Amount: summary.Total.Amount, Currency: shared.CurrencyCLP, IsApproximate: summary.Total.IsApproximate, IsIncomplete: summary.Total.IsIncomplete},
		Charged:        summary.ChargedThisMonth,
		Remaining:      summary.RemainingThisMonth,
		Rises:          risesOf(domain.RisesWithImpact(summary)),
		PendingActions: pendingActionsOf(domain.PendingActions(summary)),
		Upcoming:       upcomingListOf(domain.UpcomingCharges(summary.Subscriptions, summary.AsOf)),
	}
}

func subscriptionsOf(summary domain.Summary) []subscriptionResponse {
	hasRise := func(key domain.GroupKey) bool {
		return slices.ContainsFunc(summary.Rises, func(rise domain.Rise) bool { return rise.Key == key })
	}
	// Every list starts empty, never nil: an empty section is [] and not null.
	lines := make([]subscriptionResponse, 0, len(summary.Subscriptions))
	for _, line := range summary.Subscriptions {
		last := lastCharge(line)
		lines = append(lines, subscriptionResponse{
			LineKey:                 line.Key,
			Name:                    domain.ParentName(line.Parent),
			Alias:                   line.Alias,
			Cadence:                 line.Cadence,
			Amount:                  moneyOf(line.Amount, line.Currency),
			InPesos:                 pesosOf(line.InPesos),
			NextCharge:              unixOrNull(line.NextCharge),
			Bank:                    last.Bank,
			Card:                    last.Card,
			CardLast4:               last.CardLast4,
			HasRise:                 hasRise(line.Key),
			MarkedForCancellationAt: unixOrNull(line.MarkedForCancellationAt),
		})
	}
	return lines
}

// lastCharge tells bank and card; a manual line has no charges and shows none.
func lastCharge(line domain.Subscription) shared.Charge {
	if len(line.Charges) == 0 {
		return shared.Charge{}
	}
	return line.Charges[len(line.Charges)-1]
}

func possibleListOf(lines []domain.Subscription) []possibleResponse {
	possible := make([]possibleResponse, 0, len(lines))
	for _, line := range lines {
		possible = append(possible, possibleResponse{
			LineKey: line.Key,
			Name:    domain.ParentName(line.Parent),
			Alias:   line.Alias,
			Amount:  moneyOf(line.Amount, line.Currency),
			Charges: seenChargesOf(line.Charges),
		})
	}
	return possible
}

func seenChargesOf(charges []shared.Charge) []seenChargeResponse {
	seen := make([]seenChargeResponse, 0, len(charges))
	for _, charge := range charges {
		seen = append(seen, seenChargeResponse{Date: charge.OccurredAt.Unix(), Amount: moneyOf(charge.Amount, charge.Currency)})
	}
	return seen
}

func risesOf(rises []domain.RiseWithImpact) []riseResponse {
	responses := make([]riseResponse, 0, len(rises))
	for _, rise := range rises {
		responses = append(responses, riseOf(rise))
	}
	return responses
}

func riseOf(rise domain.RiseWithImpact) riseResponse {
	currency := rise.Rise.Currency
	return riseResponse{
		LineKey:       rise.Rise.Key,
		Before:        moneyOf(rise.Rise.Before, currency),
		After:         moneyOf(rise.Rise.After, currency),
		Since:         rise.Rise.Since.Unix(),
		MonthlyImpact: moneyOf(rise.Impact.Monthly, currency),
		YearlyImpact:  moneyOf(rise.Impact.Yearly, currency),
		IsPlanChange:  rise.Rise.IsPlanChange,
	}
}

func pendingActionsOf(actions []domain.PendingAction) []pendingActionResponse {
	responses := make([]pendingActionResponse, 0, len(actions))
	for _, action := range actions {
		responses = append(responses, pendingActionResponse{Kind: action.Kind, LineKey: action.Key})
	}
	return responses
}

func upcomingListOf(lines []domain.Subscription) []upcomingResponse {
	upcoming := make([]upcomingResponse, 0, len(lines))
	for _, line := range lines {
		upcoming = append(upcoming, upcomingResponse{LineKey: line.Key, Name: domain.ParentName(line.Parent), Alias: line.Alias, Amount: moneyOf(line.Amount, line.Currency), Date: line.NextCharge.Unix()})
	}
	return upcoming
}

func moneyOf(amount shared.Money, currency shared.Currency) moneyResponse {
	return moneyResponse{Amount: amount, Currency: currency}
}

func pesosOf(inPesos *domain.PesosAmount) *pesosResponse {
	if inPesos == nil {
		return nil
	}
	return &pesosResponse{Amount: inPesos.Amount, IsApproximate: inPesos.IsApproximate}
}

// unixOrNull answers null for a date that is not set (zero).
func unixOrNull(at time.Time) *int64 {
	if at.IsZero() {
		return nil
	}
	seconds := at.Unix()
	return &seconds
}
