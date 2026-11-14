// Package http is the subscriptions feature's HTTP edge.
package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

// Handlers are the subscriptions routes; every one of them needs a session.
type Handlers struct {
	Summary             usecases.GetSummary
	Decide              usecases.DecidePossibleSubscription
	UndoDecision        usecases.UndoDecision
	NameLine            usecases.NameLine
	SaveOverride        usecases.SaveLineOverride
	CreateManual        usecases.CreateManualSubscription
	EditManual          usecases.EditManualSubscription
	RemoveManual        usecases.RemoveManualSubscription
	KeepRise            usecases.KeepPriceRise
	CancelGuide         usecases.GetCancelGuide
	MarkForCancellation usecases.MarkForCancellation
	StopFollowing       usecases.StopFollowingLine
	Responder           httpserver.Responder
}

// MountProtected adds the routes that go behind the session.
func (h Handlers) MountProtected(router chi.Router) {
	router.Get("/api/summary", h.serve(h.summary))
	router.Post("/api/subscriptions/decisions", h.serve(h.decide))
	router.Delete("/api/subscriptions/decisions", h.serve(h.undoDecision))
	router.Put("/api/subscriptions/aliases", h.serve(h.nameLine))
	router.Put("/api/subscriptions/overrides", h.serve(h.saveOverride))
	router.Post("/api/subscriptions/manual", h.serve(h.createManual))
	router.Put("/api/subscriptions/manual/{publicId}", h.serve(h.editManual))
	router.Delete("/api/subscriptions/manual/{publicId}", h.serve(h.removeManual))
	router.Post("/api/subscriptions/price-alerts/decision", h.serve(h.keepRise))
	router.Get("/api/subscriptions/cancel-guide", h.serve(h.cancelGuide))
	router.Post("/api/subscriptions/cancellations", h.serve(h.markForCancellation))
	router.Delete("/api/subscriptions/line", h.serve(h.stopFollowing))
}

// userRequest is one request of the signed-in user.
type userRequest struct {
	exchange httpserver.Exchange
	userID   int64
}

func (u userRequest) request() *http.Request { return u.exchange.Request }

// userAction answers a request of the signed-in user; a reply without body is sent bare.
type userAction func(userRequest) (httpserver.Reply, error)

var noContent = httpserver.Reply{Status: http.StatusNoContent}

// serve reads the user from the session, runs the action and writes its reply or error.
func (h Handlers) serve(action userAction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchange := httpserver.Exchange{Writer: w, Request: r}
		reply, err := runAsUser(exchange, action)
		if err != nil {
			h.Responder.WriteError(exchange, err)
			return
		}
		if reply.Body == nil {
			w.WriteHeader(reply.Status)
			return
		}
		h.Responder.Respond(exchange, reply)
	}
}

func runAsUser(exchange httpserver.Exchange, action userAction) (httpserver.Reply, error) {
	userID, err := httpserver.UserIDFrom(exchange.Request.Context())
	if err != nil {
		return httpserver.Reply{}, err
	}
	return action(userRequest{exchange: exchange, userID: userID})
}

func (h Handlers) summary(u userRequest) (httpserver.Reply, error) {
	months, err := parseMonths(u.request().URL.Query().Get("months"))
	if err != nil {
		return httpserver.Reply{}, err
	}
	summary, err := h.Summary.Execute(u.request().Context(), u.userID, months)
	if err != nil {
		return httpserver.Reply{}, err
	}
	return httpserver.Reply{Status: http.StatusOK, Body: summaryResponseOf(summary)}, nil
}

// parseMonths is the period at the edge: missing is the default, anything that
// is not 1 to 24 is a field error.
func parseMonths(raw string) (int, error) {
	if raw == "" {
		return domain.DefaultPeriodMonths, nil
	}
	months, err := strconv.Atoi(raw)
	if err == nil {
		err = domain.ValidatePeriodMonths(months)
	}
	if err != nil {
		return 0, shared.InvalidFields(map[string]shared.MessageKey{"months": domain.KeyMonthsOutOfRange})
	}
	return months, nil
}
