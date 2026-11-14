package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/export"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
)

// Handlers are the auth routes. Register, login and resend answer 204 the same
// whether the email has an account or not.
type Handlers struct {
	Register     usecases.RegisterUser
	Login        usecases.Login
	OpenSession  usecases.OpenSession
	CloseSession usecases.CloseSession
	GetMe        usecases.GetMe
	DeleteUser   usecases.DeleteUser
	Export       export.UserCopies
	Cookies      SessionCookies
	Responder    httpserver.Responder
}

// MountPublic adds the routes that need no session. Resend is the same request
// as login: a new code replaces the previous one.
func (h Handlers) MountPublic(router chi.Router) {
	router.Post("/api/auth/register", h.register)
	router.Post("/api/auth/login", h.requestCode)
	router.Post("/api/auth/code/resend", h.requestCode)
	router.Post("/api/auth/code/verify", h.verifyCode)
	router.Post("/api/auth/logout", h.logout)
}

// MountProtected adds the routes that go behind RequireSession.
func (h Handlers) MountProtected(router chi.Router) {
	router.Get("/api/me", h.me)
	router.Delete("/api/me", h.deleteMe)
	router.Get("/api/export", h.exportCopy)
}

type registerRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	AcceptedTerms bool   `json:"acceptedTerms"`
}

type emailRequest struct {
	Email string `json:"email"`
}

type verifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type meResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
}

func (h Handlers) register(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	var request registerRequest
	if err := httpserver.DecodeJSON(exchange, &request); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	form := domain.RegistrationForm{Name: request.Name, Email: request.Email, AcceptedTerms: request.AcceptedTerms}
	if err := validateRegistration(form); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	if _, err := h.Register.Execute(r.Context(), form); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateRegistration(form domain.RegistrationForm) error {
	return shared.FieldErrors(map[string]error{
		"name":          fieldError(domain.ParseName(form.Name)),
		"email":         fieldError(domain.ParseEmail(form.Email)),
		"acceptedTerms": domain.RequireTerms(form.AcceptedTerms),
	})
}

// fieldError keeps only the error of a parse step, for FieldErrors.
func fieldError(_ string, err error) error { return err }

func (h Handlers) requestCode(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	var request emailRequest
	if err := httpserver.DecodeJSON(exchange, &request); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	if err := shared.FieldErrors(map[string]error{"email": fieldError(domain.ParseEmail(request.Email))}); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	if err := h.Login.Execute(r.Context(), request.Email); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// verifyCode checks the shape of the code first: a malformed one never uses up an attempt.
func (h Handlers) verifyCode(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	var request verifyRequest
	if err := httpserver.DecodeJSON(exchange, &request); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	code, codeErr := domain.ParseCode(request.Code)
	if err := shared.FieldErrors(map[string]error{"email": fieldError(domain.ParseEmail(request.Email)), "code": codeErr}); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	token, err := h.OpenSession.Execute(r.Context(), usecases.CodeEntry{Email: request.Email, Code: code})
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	h.Cookies.Set(w, token)
	w.WriteHeader(http.StatusNoContent)
}

// logout always answers 204 and clears the cookie, even when it no longer opened anything.
func (h Handlers) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(SessionCookieName); err == nil {
		if err := h.CloseSession.Execute(r.Context(), cookie.Value); err != nil {
			h.Responder.WriteError(httpserver.Exchange{Writer: w, Request: r}, err)
			return
		}
	}
	h.Cookies.Clear(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h Handlers) me(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	userID, err := httpserver.UserIDFrom(r.Context())
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	user, err := h.GetMe.Execute(r.Context(), userID)
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	h.Responder.Respond(exchange, httpserver.Reply{Status: http.StatusOK, Body: meResponse{ID: user.PublicID, Name: user.Name, Email: user.Email, Initials: domain.Initials(user.Name)}})
}

// deleteMe removes the account and all its data; its sessions go with it, so
// the cookie is cleared too.
func (h Handlers) deleteMe(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	userID, err := httpserver.UserIDFrom(r.Context())
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	if err := h.DeleteUser.Execute(r.Context(), userID); err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	h.Cookies.Clear(w)
	w.WriteHeader(http.StatusNoContent)
}
