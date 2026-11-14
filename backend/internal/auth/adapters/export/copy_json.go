package export

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	statementsdomain "github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	subscriptionsdomain "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// The JSON copy follows the API: ids are public ids, dates unix seconds and
// amounts {amount, currency} in hundredths. Internal ids never leave.
type copyDocument struct {
	Profile    profileCopy     `json:"profile"`
	Statements []statementCopy `json:"statements"`
	Charges    []chargeCopy    `json:"charges"`
	Decisions  []decisionCopy  `json:"decisions"`
	Aliases    []aliasCopy     `json:"aliases"`
	Manuals    []manualCopy    `json:"manuals"`
	Overrides  []overrideCopy  `json:"overrides"`
}

type profileCopy struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	AcceptedTermsAt int64  `json:"acceptedTermsAt"`
}

type statementCopy struct {
	ID          string `json:"id"`
	Bank        string `json:"bank"`
	Filename    string `json:"filename"`
	BillingDate int64  `json:"billingDate"`
	ImportedAt  int64  `json:"importedAt"`
}

type moneyCopy struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type chargeCopy struct {
	StatementID *string   `json:"statementId"`
	Bank        string    `json:"bank"`
	Source      string    `json:"source"`
	Date        int64     `json:"date"`
	Amount      moneyCopy `json:"amount"`
	Merchant    string    `json:"merchant"`
	Card        string    `json:"card"`
	CardLast4   string    `json:"cardLast4"`
}

type decisionCopy struct {
	LineKey  string `json:"lineKey"`
	Decision string `json:"decision"`
	Cadence  string `json:"cadence,omitempty"`
}

type aliasCopy struct {
	LineKey string `json:"lineKey"`
	Alias   string `json:"alias"`
}

type manualCopy struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Amount     moneyCopy `json:"amount"`
	BillingDay int       `json:"billingDay"`
}

type overrideCopy struct {
	LineKey    string `json:"lineKey"`
	Amount     *int64 `json:"amount"`
	BillingDay *int   `json:"billingDay"`
}

// EncodeJSON writes the whole copy; empty lists are [] and keyed entries go in key order.
func EncodeJSON(userCopy UserCopy) ([]byte, error) {
	body, err := json.MarshalIndent(documentOf(userCopy), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode json copy: %w", err)
	}
	return body, nil
}

func documentOf(userCopy UserCopy) copyDocument {
	return copyDocument{
		Profile:    profileCopy{ID: userCopy.Profile.PublicID, Name: userCopy.Profile.Name, Email: userCopy.Profile.Email, AcceptedTermsAt: userCopy.Profile.AcceptedTermsAt.Unix()},
		Statements: statementCopies(userCopy.Statements),
		Charges:    chargeCopies(userCopy.Saved.Charges, publicIDsOf(userCopy.Statements)),
		Decisions:  decisionCopies(userCopy.Saved.Decisions),
		Aliases:    aliasCopies(userCopy.Saved.Aliases),
		Manuals:    manualCopies(userCopy.Saved.Manuals),
		Overrides:  overrideCopies(userCopy.Saved.Overrides),
	}
}

func moneyOf(amount shared.Money, currency shared.Currency) moneyCopy {
	return moneyCopy{Amount: int64(amount), Currency: string(currency)}
}

func statementCopies(statements []statementsdomain.Statement) []statementCopy {
	copies := make([]statementCopy, 0, len(statements))
	for _, statement := range statements {
		copies = append(copies, statementCopy{ID: statement.PublicID, Bank: string(statement.Bank), Filename: statement.Filename,
			BillingDate: statement.BillingDate.Unix(), ImportedAt: statement.ImportedAt.Unix()})
	}
	return copies
}

// publicIDsOf lets a charge name its statement by the id the user knows.
func publicIDsOf(statements []statementsdomain.Statement) map[int64]string {
	ids := make(map[int64]string, len(statements))
	for _, statement := range statements {
		ids[statement.ID] = statement.PublicID
	}
	return ids
}

func chargeCopies(charges []shared.Charge, statementIDs map[int64]string) []chargeCopy {
	copies := make([]chargeCopy, 0, len(charges))
	for _, charge := range charges {
		copies = append(copies, chargeCopy{StatementID: statementIDOf(charge, statementIDs), Bank: string(charge.Bank), Source: string(charge.Source),
			Date: charge.OccurredAt.Unix(), Amount: moneyOf(charge.Amount, charge.Currency), Merchant: charge.Merchant, Card: string(charge.Card), CardLast4: charge.CardLast4})
	}
	return copies
}

func statementIDOf(charge shared.Charge, statementIDs map[int64]string) *string {
	if charge.StatementID == nil {
		return nil
	}
	publicID, isKnown := statementIDs[*charge.StatementID]
	if !isKnown {
		return nil
	}
	return &publicID
}

func decisionCopies(decisions map[subscriptionsdomain.GroupKey]subscriptionsdomain.Decision) []decisionCopy {
	copies := make([]decisionCopy, 0, len(decisions))
	for _, key := range slices.Sorted(maps.Keys(decisions)) {
		copies = append(copies, decisionCopy{LineKey: string(key), Decision: string(decisions[key].Kind), Cadence: string(decisions[key].Cadence)})
	}
	return copies
}

func aliasCopies(aliases map[subscriptionsdomain.GroupKey]string) []aliasCopy {
	copies := make([]aliasCopy, 0, len(aliases))
	for _, key := range slices.Sorted(maps.Keys(aliases)) {
		copies = append(copies, aliasCopy{LineKey: string(key), Alias: aliases[key]})
	}
	return copies
}

func manualCopies(manuals []subscriptionsdomain.ManualSubscription) []manualCopy {
	copies := make([]manualCopy, 0, len(manuals))
	for _, manual := range manuals {
		copies = append(copies, manualCopy{ID: manual.PublicID, Name: manual.Name, Amount: moneyOf(manual.Amount, manual.Currency), BillingDay: manual.BillingDay})
	}
	return copies
}

func overrideCopies(overrides map[subscriptionsdomain.GroupKey]subscriptionsdomain.LineOverride) []overrideCopy {
	copies := make([]overrideCopy, 0, len(overrides))
	for _, key := range slices.Sorted(maps.Keys(overrides)) {
		copies = append(copies, overrideCopyOf(key, overrides[key]))
	}
	return copies
}

// overrideCopyOf writes a part that is not replaced as null, like the API takes it.
func overrideCopyOf(key subscriptionsdomain.GroupKey, override subscriptionsdomain.LineOverride) overrideCopy {
	copied := overrideCopy{LineKey: string(key)}
	if override.Amount > 0 {
		amount := int64(override.Amount)
		copied.Amount = &amount
	}
	if override.BillingDay > 0 {
		day := override.BillingDay
		copied.BillingDay = &day
	}
	return copied
}
