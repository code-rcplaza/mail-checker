package http

import (
	"net/http"
	"strings"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

const (
	lineKeyField  = "lineKey"
	groupKeyField = "groupKey"
)

type decisionRequest struct {
	GroupKey string `json:"groupKey"`
	Decision string `json:"decision"`
	Cadence  string `json:"cadence"`
}

type aliasRequest struct {
	LineKey string `json:"lineKey"`
	Alias   string `json:"alias"`
}

type overrideRequest struct {
	LineKey    string        `json:"lineKey"`
	Amount     *shared.Money `json:"amount"`
	BillingDay *int          `json:"billingDay"`
}

type lineKeyRequest struct {
	LineKey string `json:"lineKey"`
}

type cancelGuideResponse struct {
	URL string `json:"url"`
}

// lineKeyInput is a key the user sent and the field it came in.
type lineKeyInput struct {
	field string
	raw   string
}

// requireKey rejects a missing or blank key at the edge, on the field it came in.
func requireKey(input lineKeyInput) (domain.GroupKey, error) {
	if strings.TrimSpace(input.raw) == "" {
		return "", shared.InvalidFields(map[string]shared.MessageKey{input.field: domain.KeyLineKeyRequired})
	}
	return domain.GroupKey(input.raw), nil
}

func queryKey(u userRequest, field string) (domain.GroupKey, error) {
	return requireKey(lineKeyInput{field: field, raw: u.request().URL.Query().Get(field)})
}

// bodyLineKey reads a {lineKey} body, the shape of every one-line action.
func bodyLineKey(u userRequest) (domain.GroupKey, error) {
	var request lineKeyRequest
	if err := httpserver.DecodeJSON(u.exchange, &request); err != nil {
		return "", err
	}
	return requireKey(lineKeyInput{field: lineKeyField, raw: request.LineKey})
}

func (h Handlers) decide(u userRequest) (httpserver.Reply, error) {
	var request decisionRequest
	if err := httpserver.DecodeJSON(u.exchange, &request); err != nil {
		return httpserver.Reply{}, err
	}
	key, err := requireKey(lineKeyInput{field: groupKeyField, raw: request.GroupKey})
	if err != nil {
		return httpserver.Reply{}, err
	}
	input := usecases.DecisionInput{Key: key, Action: request.Decision, Cadence: request.Cadence}
	return noContent, h.Decide.Execute(u.request().Context(), u.userID, input)
}

func (h Handlers) undoDecision(u userRequest) (httpserver.Reply, error) {
	key, err := queryKey(u, groupKeyField)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return noContent, h.UndoDecision.Execute(u.request().Context(), u.userID, key)
}

func (h Handlers) nameLine(u userRequest) (httpserver.Reply, error) {
	var request aliasRequest
	if err := httpserver.DecodeJSON(u.exchange, &request); err != nil {
		return httpserver.Reply{}, err
	}
	key, err := requireKey(lineKeyInput{field: lineKeyField, raw: request.LineKey})
	if err != nil {
		return httpserver.Reply{}, err
	}
	return noContent, h.NameLine.Execute(u.request().Context(), u.userID, domain.LineAlias{Key: key, Alias: request.Alias})
}

func (h Handlers) saveOverride(u userRequest) (httpserver.Reply, error) {
	var request overrideRequest
	if err := httpserver.DecodeJSON(u.exchange, &request); err != nil {
		return httpserver.Reply{}, err
	}
	key, err := requireKey(lineKeyInput{field: lineKeyField, raw: request.LineKey})
	if err != nil {
		return httpserver.Reply{}, err
	}
	if err := domain.ValidateOverrideFields(request.Amount, request.BillingDay); err != nil {
		return httpserver.Reply{}, err
	}
	input := usecases.LineOverrideInput{Key: key, Amount: request.Amount, BillingDay: request.BillingDay}
	return noContent, h.SaveOverride.Execute(u.request().Context(), u.userID, input)
}

func (h Handlers) keepRise(u userRequest) (httpserver.Reply, error) {
	key, err := bodyLineKey(u)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return noContent, h.KeepRise.Execute(u.request().Context(), u.userID, key)
}

func (h Handlers) cancelGuide(u userRequest) (httpserver.Reply, error) {
	key, err := queryKey(u, lineKeyField)
	if err != nil {
		return httpserver.Reply{}, err
	}
	link, err := h.CancelGuide.Execute(u.request().Context(), u.userID, key)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return httpserver.Reply{Status: http.StatusOK, Body: cancelGuideResponse{URL: link}}, nil
}

func (h Handlers) markForCancellation(u userRequest) (httpserver.Reply, error) {
	key, err := bodyLineKey(u)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return noContent, h.MarkForCancellation.Execute(u.request().Context(), u.userID, key)
}

func (h Handlers) stopFollowing(u userRequest) (httpserver.Reply, error) {
	key, err := queryKey(u, lineKeyField)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return noContent, h.StopFollowing.Execute(u.request().Context(), u.userID, key)
}
