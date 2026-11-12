package domain

// Service is a known merchant from the catalog. An umbrella service (an app
// store, for instance) bills many unrelated subscriptions under one name.
type Service struct {
	Name       string
	IsUmbrella bool
}
