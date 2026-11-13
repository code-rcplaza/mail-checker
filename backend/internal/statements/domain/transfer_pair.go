package domain

import (
	"time"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// TransferPair is the same debt in pesos and in dollars: the bank moves the
// international balance to the national statement, which gives the rate it
// really charged.
type TransferPair struct {
	Pesos   shared.Money
	Dollars shared.Money
	Date    time.Time
}
