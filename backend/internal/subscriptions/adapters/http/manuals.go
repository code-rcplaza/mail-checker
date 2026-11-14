package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

type manualRequest struct {
	Name       string       `json:"name"`
	Amount     shared.Money `json:"amount"`
	Currency   string       `json:"currency"`
	BillingDay int          `json:"billingDay"`
}

type createdResponse struct {
	ID string `json:"id"`
}

// readManual decodes the form and checks every field at once, so the front
// shows each error next to its field.
func readManual(u userRequest) (domain.ManualSubscriptionInput, error) {
	var request manualRequest
	if err := httpserver.DecodeJSON(u.exchange, &request); err != nil {
		return domain.ManualSubscriptionInput{}, err
	}
	input := domain.ManualSubscriptionInput{Name: request.Name, Amount: request.Amount, Currency: request.Currency, BillingDay: request.BillingDay}
	if err := domain.ValidateManualFields(input); err != nil {
		return domain.ManualSubscriptionInput{}, err
	}
	return input, nil
}

// manualPublicID is the id in the path; a malformed one matches no manual and
// answers 404 like one of another user.
func manualPublicID(u userRequest) string { return chi.URLParam(u.request(), "publicId") }

func (h Handlers) createManual(u userRequest) (httpserver.Reply, error) {
	input, err := readManual(u)
	if err != nil {
		return httpserver.Reply{}, err
	}
	manual, err := h.CreateManual.Execute(u.request().Context(), u.userID, input)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return httpserver.Reply{Status: http.StatusCreated, Body: createdResponse{ID: manual.PublicID}}, nil
}

func (h Handlers) editManual(u userRequest) (httpserver.Reply, error) {
	input, err := readManual(u)
	if err != nil {
		return httpserver.Reply{}, err
	}
	edit := usecases.ManualSubscriptionEdit{PublicID: manualPublicID(u), Input: input}
	return noContent, h.EditManual.Execute(u.request().Context(), u.userID, edit)
}

func (h Handlers) removeManual(u userRequest) (httpserver.Reply, error) {
	return noContent, h.RemoveManual.Execute(u.request().Context(), u.userID, manualPublicID(u))
}
