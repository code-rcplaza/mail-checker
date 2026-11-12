package domain

import (
	"reflect"
	"strings"
	"testing"
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// keyedLine builds a line charged on the given days of consecutive months from July 2026.
func keyedLine(last4 string, days ...int) []shared.Charge {
	line := make([]shared.Charge, 0, len(days))
	for i, day := range days {
		c := lineCharge("c", 1299000, shared.CurrencyCLP, time.July+time.Month(i), day)
		c.CardLast4 = last4
		line = append(line, c)
	}
	return line
}

func TestLineKeys_ParentCardAndMedianDay(t *testing.T) {
	got := LineKeys("service:Netflix", [][]shared.Charge{keyedLine("4821", 1, 2, 1)})

	if want := []GroupKey{"service:Netflix|4821|d01"}; !reflect.DeepEqual(got, want) {
		t.Errorf("keys = %v, want %v", got, want)
	}
}

func TestLineKeys_LineWithoutCardLeavesItEmpty(t *testing.T) {
	got := LineKeys("service:Netflix", [][]shared.Charge{keyedLine("", 1, 1)})

	if want := []GroupKey{"service:Netflix||d01"}; !reflect.DeepEqual(got, want) {
		t.Errorf("keys = %v, want %v", got, want)
	}
}

func TestLineKeys_CollidingLinesNumberedByFirstCharge(t *testing.T) {
	later := keyedLine("4821", 1, 1)
	for i := range later {
		later[i].OccurredAt = later[i].OccurredAt.AddDate(0, 1, 0)
	}
	earlier := keyedLine("4821", 1, 1, 1)

	got := LineKeys("umbrella:App Store", [][]shared.Charge{later, earlier})

	want := []GroupKey{"umbrella:App Store|4821|d01#2", "umbrella:App Store|4821|d01"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("keys = %v, want %v", got, want)
	}
}

func TestLineKeys_SameKeyWithAShorterPeriod(t *testing.T) {
	sixMonths := keyedLine("4821", 1, 2, 1, 1, 2, 1)
	threeMonths := sixMonths[3:]

	long := LineKeys("service:Netflix", [][]shared.Charge{sixMonths})
	short := LineKeys("service:Netflix", [][]shared.Charge{threeMonths})

	want := []GroupKey{"service:Netflix|4821|d01"}
	if !reflect.DeepEqual(long, want) || !reflect.DeepEqual(short, want) {
		t.Errorf("6 months = %v, 3 months = %v; want both %v", long, short, want)
	}
}

func TestLineKeys_UnknownMerchantKeepsItsParent(t *testing.T) {
	got := LineKeys("merchant:SEGUROS CLC S A:CLP", [][]shared.Charge{keyedLine("4821", 4, 6)})

	if want := []GroupKey{"merchant:SEGUROS CLC S A:CLP|4821|d04"}; !reflect.DeepEqual(got, want) {
		t.Errorf("keys = %v, want %v", got, want)
	}
}

func TestLineKeys_NoInternalIDsOrEmails(t *testing.T) {
	line := keyedLine("4821", 1, 1)
	accountID := int64(987654)
	for i := range line {
		line[i].AccountID = &accountID
		line[i].ExternalID = "<notice-1@example.com>"
	}

	key := string(LineKeys("service:Netflix", [][]shared.Charge{line})[0])

	if key == "" || strings.Contains(key, "987654") || strings.Contains(key, "@") {
		t.Errorf("key = %q, want a non-empty key without account ids or emails", key)
	}
}
