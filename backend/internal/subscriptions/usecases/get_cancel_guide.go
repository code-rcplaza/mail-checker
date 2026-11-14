package usecases

import (
	"context"

	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
)

// GetCancelGuide gives the link where the user cancels a line's service. A
// merchant the catalog does not know, or a service whose link is not verified
// yet, has none: the front shows the generic steps.
type GetCancelGuide struct {
	lines ChargeLines
}

func NewGetCancelGuide(lines ChargeLines) GetCancelGuide { return GetCancelGuide{lines: lines} }

func (u GetCancelGuide) Execute(ctx context.Context, userID int64, key domain.GroupKey) (string, error) {
	line, err := u.lines.findLine(ctx, userID, key)
	if err != nil {
		return "", err
	}
	return u.cancelLinkOf(line), nil
}

// cancelLinkOf asks the catalog about the line's last charge: the merchant it is billed as now.
func (u GetCancelGuide) cancelLinkOf(line domain.ChargeGroup) string {
	if len(line.Charges) == 0 {
		return ""
	}
	service, isKnown := u.lines.Catalog.Identify(line.Charges[len(line.Charges)-1].Merchant)
	if !isKnown {
		return ""
	}
	return service.CancelURL
}
