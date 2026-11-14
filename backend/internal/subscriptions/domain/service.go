package domain

import (
	"net/url"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

// Service is a known merchant from the catalog. An umbrella service (an app
// store, for instance) bills many unrelated subscriptions under one name.
type Service struct {
	Name       string
	IsUmbrella bool
	CancelURL  string // where to cancel or manage it; empty when not verified
}

// NewService rejects a cancel link that is not https, so an http or odd
// scheme link never leaves through the API.
func NewService(candidate Service) (Service, error) {
	if err := validateCancelURL(candidate.CancelURL); err != nil {
		return Service{}, err
	}
	return candidate, nil
}

func validateCancelURL(cancelURL string) error {
	if cancelURL == "" || isHTTPSLink(cancelURL) {
		return nil
	}
	return shared.Invalid(KeyCancelURLInvalid, nil)
}

func isHTTPSLink(link string) bool {
	parsed, err := url.Parse(link)
	if err != nil {
		return false
	}
	return parsed.Scheme == "https" && parsed.Host != ""
}
