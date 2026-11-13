package usecases

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
)

func lineOn(number int, date time.Time) domain.StatementLine {
	line := statementLine(number, "STREAMING CO")
	line.Charge.Date = date
	return line
}

func TestImportStatement_ReturnsTheMonthsRead(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.lines = []domain.StatementLine{
		lineOn(2, time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)),
		lineOn(3, time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)),
		lineOn(4, time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)),
	}

	got, err := fixture.useCase.Execute(context.Background(), userOne, bciInput())

	if err != nil {
		t.Fatalf("Execute error = %v, want nil", err)
	}
	want := []domain.Month{{Year: 2026, Month: time.September}, {Year: 2026, Month: time.October}}
	if !reflect.DeepEqual(got.Months, want) {
		t.Errorf("months = %v, want %v", got.Months, want)
	}
}

func TestImportStatement_RejectedImportHasNoResult(t *testing.T) {
	fixture := newImportFixture()
	fixture.bci.lines = nil

	got, err := fixture.useCase.Execute(context.Background(), userOne, bciInput())

	if err == nil {
		t.Fatalf("Execute error = nil, want a rejection")
	}
	if !reflect.DeepEqual(got, ImportResult{}) {
		t.Errorf("result = %+v, want empty", got)
	}
}
