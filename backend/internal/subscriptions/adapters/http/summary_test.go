package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

var asOf = time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)

func statementID(id int64) *int64 { return &id }

// sampleSummary has one active line in dollars with a rise and a mark, one possible, and a manual.
func sampleSummary() domain.Summary {
	rate := shared.ExchangeRate{PesosPerDollar: 95000, Date: asOf}
	netflix := domain.Subscription{
		Key: "service:Netflix||d20", Parent: "service:Netflix", Status: domain.StatusActive, Cadence: domain.CadenceMonthly,
		Amount: 2000, Currency: shared.CurrencyUSD, InPesos: &domain.PesosAmount{Amount: 1900000, IsApproximate: true},
		NextCharge: time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC), MarkedForCancellationAt: asOf.AddDate(0, 0, -1),
		Charges: []shared.Charge{{ID: 99, UserID: 7, StatementID: statementID(5), Bank: "bci", ExternalID: "st:1", Merchant: "NETFLIX", Card: "credito", CardLast4: "1234", Amount: 2000, Currency: shared.CurrencyUSD, OccurredAt: asOf.AddDate(0, 0, -24)}},
	}
	gym := domain.Subscription{Key: "manual:m1", Parent: "manual:m1", Alias: "Gimnasio", Status: domain.StatusActive, Cadence: domain.CadenceMonthly,
		Amount: 2500000, Currency: shared.CurrencyCLP, InPesos: &domain.PesosAmount{Amount: 2500000}, NextCharge: time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)}
	music := domain.Subscription{Key: "merchant:MUSIC APP:CLP||d05", Parent: "merchant:MUSIC APP:CLP", Status: domain.StatusPossible, Amount: 399000, Currency: shared.CurrencyCLP,
		Charges: []shared.Charge{{ID: 42, Merchant: "MUSIC APP", Amount: 399000, Currency: shared.CurrencyCLP, OccurredAt: asOf.AddDate(0, 0, -8)}}}
	return domain.Summary{
		Subscriptions:      []domain.Subscription{netflix, gym},
		Possible:           []domain.Subscription{music},
		Inactive:           []domain.Subscription{{Key: "merchant:OLD:CLP||d01"}},
		OtherCharges:       []shared.Charge{{Merchant: "SUPERMERCADO"}},
		ChargedThisMonth:   3,
		RemainingThisMonth: 5,
		Rises:              []domain.Rise{{Key: netflix.Key, Before: 1800, After: 2000, Currency: shared.CurrencyUSD, Since: asOf.AddDate(0, 0, -24)}},
		Total:              domain.Total{Amount: 4400000, IsApproximate: true},
		ExchangeRate:       &rate,
		AsOf:               asOf,
	}
}

func encoded(t *testing.T, summary domain.Summary) (string, map[string]any) {
	t.Helper()
	body, err := json.Marshal(summaryResponseOf(summary))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return string(body), decoded
}

func TestSummaryResponse_Subscriptions(t *testing.T) {
	body, _ := encoded(t, sampleSummary())
	nextCharge := fmt.Sprintf(`"nextCharge":%d`, time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC).Unix())
	markedAt := fmt.Sprintf(`"markedForCancellationAt":%d`, asOf.AddDate(0, 0, -1).Unix())

	for _, want := range []string{
		`"lineKey":"service:Netflix||d20"`, `"name":"Netflix"`, `"cadence":"monthly"`, nextCharge,
		`"bank":"bci"`, `"card":"credito"`, `"cardLast4":"1234"`, `"hasRise":true`, markedAt,
		`"amount":{"amount":2000,"currency":"USD"}`, `"inPesos":{"amount":1900000,"isApproximate":true}`,
		`"alias":"Gimnasio"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("summary JSON lacks %s:\n%s", want, body)
		}
	}
}

func TestSummaryResponse_HasNoInternalsNorHiddenSections(t *testing.T) {
	body, decoded := encoded(t, sampleSummary())

	for _, internal := range []string{`"id":99`, `"userId"`, `"statementId"`, `"externalId"`, "st:1", `"ID"`, "SUPERMERCADO", "OLD:CLP", "inactive", "otherCharges"} {
		if strings.Contains(body, internal) {
			t.Errorf("summary JSON contains %s", internal)
		}
	}
	if _, isThere := decoded["exchangeRate"]; isThere {
		t.Errorf("summary JSON has the exchange rate, which is not in the contract")
	}
}

func TestSummaryResponse_Possible(t *testing.T) {
	_, decoded := encoded(t, sampleSummary())

	possible := decoded["possible"].([]any)[0].(map[string]any)
	if possible["lineKey"] != "merchant:MUSIC APP:CLP||d05" || possible["name"] != "MUSIC APP" {
		t.Errorf("possible = %v, want its line key and name", possible)
	}
	charges := possible["charges"].([]any)
	if len(charges) != 1 || charges[0].(map[string]any)["date"] != float64(asOf.AddDate(0, 0, -8).Unix()) {
		t.Errorf("possible charges = %v, want each seen charge with date and amount", charges)
	}
	for _, absent := range []string{"nextCharge", "cadence"} {
		if _, isThere := possible[absent]; isThere {
			t.Errorf("possible has %s, want none while undecided", absent)
		}
	}
}

func TestSummaryResponse_Total(t *testing.T) {
	cases := []struct {
		name  string
		total domain.Total
		want  string
	}{
		{"converted dollars are approximate", domain.Total{Amount: 4400000, IsApproximate: true}, `"total":{"amount":4400000,"currency":"CLP","isApproximate":true,"isIncomplete":false}`},
		{"no rate is incomplete", domain.Total{Amount: 2500000, IsIncomplete: true}, `"total":{"amount":2500000,"currency":"CLP","isApproximate":false,"isIncomplete":true}`},
		{"only pesos has no mark", domain.Total{Amount: 2500000}, `"total":{"amount":2500000,"currency":"CLP","isApproximate":false,"isIncomplete":false}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			summary := sampleSummary()
			summary.Total = tc.total

			if body, _ := encoded(t, summary); !strings.Contains(body, tc.want) {
				t.Errorf("summary JSON lacks %s:\n%s", tc.want, body)
			}
		})
	}
}

func TestSummaryResponse_MonthCountsRisesActionsAndUpcoming(t *testing.T) {
	body, decoded := encoded(t, sampleSummary())

	for _, want := range []string{`"charged":3`, `"remaining":5`,
		`"rises":[{"lineKey":"service:Netflix||d20","before":{"amount":1800,"currency":"USD"},"after":{"amount":2000,"currency":"USD"},"since":` + fmt.Sprint(asOf.AddDate(0, 0, -24).Unix()) + `,"monthlyImpact":{"amount":200,"currency":"USD"},"yearlyImpact":{"amount":2400,"currency":"USD"},"isPlanChange":false}]`,
		`"pendingActions":[{"kind":"rise","lineKey":"service:Netflix||d20"},{"kind":"possible","lineKey":"merchant:MUSIC APP:CLP||d05"}]`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("summary JSON lacks %s:\n%s", want, body)
		}
	}
	upcoming := decoded["upcoming"].([]any)
	if len(upcoming) != 2 || upcoming[0].(map[string]any)["alias"] != "Gimnasio" || upcoming[1].(map[string]any)["name"] != "Netflix" {
		t.Errorf("upcoming = %v, want the manual (15th) before Netflix (20th)", upcoming)
	}
}

func TestSummaryResponse_EmptyListsAreNotNull(t *testing.T) {
	body, _ := encoded(t, domain.Summary{AsOf: asOf})

	for _, want := range []string{`"subscriptions":[]`, `"possible":[]`, `"rises":[]`, `"pendingActions":[]`, `"upcoming":[]`, `"charged":0`, `"remaining":0`, `"total":{"amount":0`} {
		if !strings.Contains(body, want) {
			t.Errorf("empty summary JSON lacks %s:\n%s", want, body)
		}
	}
}

// summaryRouter serves the summary of a user without any statement over a real database.
func summaryRouter(t *testing.T) http.Handler {
	t.Helper()
	return newTestServer(t).router
}

func TestGetSummary_MonthsParameter(t *testing.T) {
	router := summaryRouter(t)
	cases := []struct {
		query string
		want  int
	}{
		{"?months=12", http.StatusOK},
		{"", http.StatusOK},
		{"?months=24", http.StatusOK},
		{"?months=0", http.StatusBadRequest},
		{"?months=-1", http.StatusBadRequest},
		{"?months=25", http.StatusBadRequest},
		{"?months=abc", http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/summary"+tc.query, nil))

			if recorder.Code != tc.want {
				t.Errorf("GET /api/summary%s = %d %s, want %d", tc.query, recorder.Code, recorder.Body.String(), tc.want)
			}
			if tc.want == http.StatusBadRequest && !strings.Contains(recorder.Body.String(), `"fields":{"months":"subscriptions.months_out_of_range"}`) {
				t.Errorf("body = %s, want the months field", recorder.Body.String())
			}
		})
	}
}

func TestGetSummary_UserWithoutStatements(t *testing.T) {
	recorder := httptest.NewRecorder()
	summaryRouter(t).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/summary", nil))

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"subscriptions":[]`) || !strings.Contains(recorder.Body.String(), `"total":{"amount":0`) {
		t.Errorf("summary = %d %s, want 200 with empty lists and a zero total", recorder.Code, recorder.Body.String())
	}
}
