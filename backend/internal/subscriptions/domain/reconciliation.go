package domain

import (
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type ChargePair struct {
	Statement shared.Charge
	Email     shared.Charge
}

// Reconciliation is the differentiated view: which statement charges were
// notified by email, which never were (the key finding), and which emails
// are not yet in an imported statement.
type Reconciliation struct {
	Notified   []ChargePair
	Unnotified []shared.Charge
	EmailOnly  []shared.Charge
}

func BuildReconciliation(charges []shared.Charge, identify IdentifyServiceFunc, normalize NormalizeMerchantFunc) Reconciliation {
	emailByStatement := pairStatementsWithEmails(charges, identify, normalize)
	pairedEmails := pairedEmailIndexes(emailByStatement)
	reconciliation := Reconciliation{
		Notified:   []ChargePair{},
		Unnotified: []shared.Charge{},
		EmailOnly:  []shared.Charge{},
	}
	for i, c := range charges {
		emailIndex, isPaired := emailByStatement[i]
		switch {
		case isPaired:
			reconciliation.Notified = append(reconciliation.Notified, ChargePair{Statement: c, Email: charges[emailIndex]})
		case c.Source == shared.SourceStatement:
			reconciliation.Unnotified = append(reconciliation.Unnotified, c)
		case !pairedEmails[i]:
			reconciliation.EmailOnly = append(reconciliation.EmailOnly, c)
		}
	}
	return reconciliation
}

func pairedEmailIndexes(emailByStatement map[int]int) map[int]bool {
	paired := make(map[int]bool, len(emailByStatement))
	for _, emailIndex := range emailByStatement {
		paired[emailIndex] = true
	}
	return paired
}
