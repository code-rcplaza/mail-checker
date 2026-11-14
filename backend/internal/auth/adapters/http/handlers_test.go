package http

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/mailer"
	authsqlite "github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
)

// authApp is the auth routes over a real database; the codes go to console.
type authApp struct {
	router  http.Handler
	console *bytes.Buffer
}

// fixedRandom fills every read with the same byte; enough for tests.
type fixedRandom byte

func (f fixedRandom) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(f)
	}
	return len(p), nil
}

func newAuthApp(t *testing.T) authApp {
	t.Helper()
	db := sqlitetest.OpenMigrated(t)
	now := func() time.Time { return testNow }
	console := &bytes.Buffer{}
	users := authsqlite.NewUsers(db, now)
	codes := authsqlite.NewLoginCodes(db)
	sessions := authsqlite.NewSessions(db)
	delivery := usecases.NewCodeDelivery(usecases.CodeDeliveryDependencies{Codes: codes, Sends: authsqlite.NewCodeSends(db), Sender: mailer.ConsoleSender{Writer: console}, Random: fixedRandom(3), Now: now})
	ids := 0
	newPublicID := func() string { ids++; return fmt.Sprintf("0193f1a2-0000-7000-8000-%012d", ids) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	responder := httpserver.Responder{Logger: logger}
	handlers := Handlers{
		Register:     usecases.NewRegisterUser(usecases.RegisterUserDependencies{Users: users, Delivery: delivery, NewPublicID: newPublicID, Now: now}),
		Login:        usecases.NewLogin(users, delivery),
		OpenSession:  usecases.NewOpenSession(usecases.OpenSessionDependencies{Verify: usecases.NewVerifyCode(usecases.VerifyCodeDependencies{Users: users, Codes: codes, Now: now}), Sessions: sessions, Random: fixedRandom(9), Now: now}),
		CloseSession: usecases.NewCloseSession(sessions),
		GetMe:        usecases.NewGetMe(users),
		DeleteUser:   usecases.NewDeleteUser(users),
		Export:       copiesOver(db),
		Cookies:      SessionCookies{},
		Responder:    responder,
	}
	router := httpserver.NewRouter(logger, httpserver.AlwaysHealthy)
	handlers.MountPublic(router)
	router.Group(func(protected chi.Router) {
		protected.Use(RequireSession{Validate: usecases.NewValidateSession(sessions, now), Responder: responder}.Handler)
		handlers.MountProtected(protected)
	})
	return authApp{router: router, console: console}
}

func (a authApp) do(method, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	recorder := httptest.NewRecorder()
	a.router.ServeHTTP(recorder, request)
	return recorder
}

var sentCode = regexp.MustCompile(`login code for (\S+): (\d{6})`)

// codesSent lists the codes written to the console, in order.
func (a authApp) codesSent() []string {
	var codes []string
	for _, match := range sentCode.FindAllStringSubmatch(a.console.String(), -1) {
		codes = append(codes, match[2])
	}
	return codes
}

func sessionCookieOf(recorder *httptest.ResponseRecorder) *http.Cookie {
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == SessionCookieName {
			return cookie
		}
	}
	return nil
}

const anaForm = `{"name":"Ana Pérez","email":"ana@example.com","acceptedTerms":true}`

// signIn registers Ana and verifies her code; it returns the session cookie.
func (a authApp) signIn(t *testing.T) *http.Cookie {
	t.Helper()
	if got := a.do(http.MethodPost, "/api/auth/register", anaForm); got.Code != http.StatusNoContent {
		t.Fatalf("register = %d %s", got.Code, got.Body.String())
	}
	codes := a.codesSent()
	verified := a.do(http.MethodPost, "/api/auth/code/verify", fmt.Sprintf(`{"email":"ana@example.com","code":"%s"}`, codes[len(codes)-1]))
	cookie := sessionCookieOf(verified)
	if verified.Code != http.StatusNoContent || cookie == nil {
		t.Fatalf("verify = %d %s, cookie %v", verified.Code, verified.Body.String(), cookie)
	}
	return cookie
}

type errorReply struct {
	Error struct {
		Key    string            `json:"key"`
		Fields map[string]string `json:"fields"`
	} `json:"error"`
}

func errorOf(t *testing.T, recorder *httptest.ResponseRecorder) errorReply {
	t.Helper()
	var reply errorReply
	if err := json.Unmarshal(recorder.Body.Bytes(), &reply); err != nil {
		t.Fatalf("body %q is not an error: %v", recorder.Body.String(), err)
	}
	return reply
}

func TestRegister(t *testing.T) {
	app := newAuthApp(t)

	first := app.do(http.MethodPost, "/api/auth/register", anaForm)
	again := app.do(http.MethodPost, "/api/auth/register", `{"name":"Otra","email":"ANA@example.com","acceptedTerms":true}`)

	if first.Code != http.StatusNoContent || again.Code != http.StatusNoContent {
		t.Errorf("register = %d, again = %d; want 204 both (neutral answer)", first.Code, again.Code)
	}
	if got := len(app.codesSent()); got != 2 {
		t.Errorf("codes sent = %d, want one per register (the owner gets the second)", got)
	}
}

func TestRegister_InvalidFields(t *testing.T) {
	app := newAuthApp(t)

	got := app.do(http.MethodPost, "/api/auth/register", `{"name":" ","email":"ana@","acceptedTerms":false}`)

	reply := errorOf(t, got)
	want := map[string]string{"name": "auth.name_required", "email": "auth.email_invalid", "acceptedTerms": "auth.terms_required"}
	if got.Code != http.StatusBadRequest || reply.Error.Key != "common.validation_failed" || fmt.Sprint(reply.Error.Fields) != fmt.Sprint(want) {
		t.Errorf("register = %d %+v, want 400 with fields %v", got.Code, reply, want)
	}
	if len(app.codesSent()) != 0 {
		t.Errorf("a code was sent for an invalid form")
	}
}

func TestRegister_UnknownFieldIsInvalidJSON(t *testing.T) {
	got := newAuthApp(t).do(http.MethodPost, "/api/auth/register", `{"name":"Ana","email":"ana@example.com","acceptedTerms":true,"admin":true}`)

	if got.Code != http.StatusBadRequest || errorOf(t, got).Error.Key != "common.invalid_json" {
		t.Errorf("register = %d %s, want 400 common.invalid_json", got.Code, got.Body.String())
	}
}

func TestLoginAndResend(t *testing.T) {
	for _, path := range []string{"/api/auth/login", "/api/auth/code/resend"} {
		t.Run(path, func(t *testing.T) {
			app := newAuthApp(t)
			app.do(http.MethodPost, "/api/auth/register", anaForm)

			registered := app.do(http.MethodPost, path, `{"email":"ana@example.com"}`)
			unknown := app.do(http.MethodPost, path, `{"email":"nadie@example.com"}`)
			invalid := app.do(http.MethodPost, path, `{"email":"nadie@"}`)

			if registered.Code != http.StatusNoContent || unknown.Code != http.StatusNoContent {
				t.Errorf("registered = %d, unknown = %d; want 204 both", registered.Code, unknown.Code)
			}
			if invalid.Code != http.StatusBadRequest || errorOf(t, invalid).Error.Fields["email"] != "auth.email_invalid" {
				t.Errorf("invalid email = %d %s, want 400 with the email field", invalid.Code, invalid.Body.String())
			}
			if got := len(app.codesSent()); got != 2 {
				t.Errorf("codes sent = %d, want register + the registered request", got)
			}
		})
	}
}

func TestLogin_OverTheLimitIs429(t *testing.T) {
	app := newAuthApp(t)
	app.do(http.MethodPost, "/api/auth/register", anaForm)
	app.do(http.MethodPost, "/api/auth/login", `{"email":"ana@example.com"}`)
	app.do(http.MethodPost, "/api/auth/code/resend", `{"email":"ana@example.com"}`)

	got := app.do(http.MethodPost, "/api/auth/login", `{"email":"ana@example.com"}`)

	if got.Code != http.StatusTooManyRequests || errorOf(t, got).Error.Key != "auth.too_many_sends" {
		t.Errorf("fourth request = %d %s, want 429 auth.too_many_sends", got.Code, got.Body.String())
	}
}

func TestVerify(t *testing.T) {
	app := newAuthApp(t)
	app.do(http.MethodPost, "/api/auth/register", anaForm)
	code := app.codesSent()[0]
	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}

	bad := app.do(http.MethodPost, "/api/auth/code/verify", `{"email":"ana@example.com","code":"`+wrong+`"}`)
	malformed := app.do(http.MethodPost, "/api/auth/code/verify", `{"email":"ana@example.com","code":"12a"}`)
	good := app.do(http.MethodPost, "/api/auth/code/verify", `{"email":"ana@example.com","code":" `+code+` "}`)

	if bad.Code != http.StatusBadRequest || errorOf(t, bad).Error.Key != "auth.code_invalid" || sessionCookieOf(bad) != nil {
		t.Errorf("wrong code = %d %s, want 400 auth.code_invalid without cookie", bad.Code, bad.Body.String())
	}
	if malformed.Code != http.StatusBadRequest || errorOf(t, malformed).Error.Fields["code"] != "auth.code_format_invalid" {
		t.Errorf("malformed code = %d %s, want 400 with the code field", malformed.Code, malformed.Body.String())
	}
	if good.Code != http.StatusNoContent || sessionCookieOf(good) == nil {
		t.Errorf("right code (with spaces) = %d %s, want 204 with the session cookie", good.Code, good.Body.String())
	}
}

func TestVerify_FifthWrongAttemptExhaustsTheCode(t *testing.T) {
	app := newAuthApp(t)
	app.do(http.MethodPost, "/api/auth/register", anaForm)
	code := app.codesSent()[0]
	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}
	for range 4 {
		app.do(http.MethodPost, "/api/auth/code/verify", `{"email":"ana@example.com","code":"`+wrong+`"}`)
	}
	malformed := app.do(http.MethodPost, "/api/auth/code/verify", `{"email":"ana@example.com","code":"1"}`)

	fifth := app.do(http.MethodPost, "/api/auth/code/verify", `{"email":"ana@example.com","code":"`+wrong+`"}`)

	if malformed.Code != http.StatusBadRequest || fifth.Code != http.StatusBadRequest || errorOf(t, fifth).Error.Key != "auth.code_attempts_exceeded" {
		t.Errorf("fifth wrong = %d %s, want 400 auth.code_attempts_exceeded (the malformed one did not count)", fifth.Code, fifth.Body.String())
	}
}

func TestMe(t *testing.T) {
	app := newAuthApp(t)
	cookie := app.signIn(t)

	got := app.do(http.MethodGet, "/api/me", "", cookie)
	anonymous := app.do(http.MethodGet, "/api/me", "")

	var me map[string]any
	if err := json.Unmarshal(got.Body.Bytes(), &me); err != nil || got.Code != http.StatusOK {
		t.Fatalf("me = %d %s", got.Code, got.Body.String())
	}
	want := map[string]any{"id": "0193f1a2-0000-7000-8000-000000000001", "name": "Ana Pérez", "email": "ana@example.com", "initials": "AP"}
	if fmt.Sprint(me) != fmt.Sprint(want) {
		t.Errorf("me = %v, want exactly %v (no internal id)", me, want)
	}
	if anonymous.Code != http.StatusUnauthorized {
		t.Errorf("me without session = %d, want 401", anonymous.Code)
	}
}

func TestLogout(t *testing.T) {
	app := newAuthApp(t)
	cookie := app.signIn(t)

	out := app.do(http.MethodPost, "/api/auth/logout", "", cookie)
	again := app.do(http.MethodPost, "/api/auth/logout", "", cookie)
	without := app.do(http.MethodPost, "/api/auth/logout", "")
	me := app.do(http.MethodGet, "/api/me", "", cookie)

	cleared := sessionCookieOf(out)
	if out.Code != http.StatusNoContent || cleared == nil || cleared.MaxAge >= 0 {
		t.Errorf("logout = %d, cookie %v; want 204 with the cookie cleared", out.Code, cleared)
	}
	if again.Code != http.StatusNoContent || without.Code != http.StatusNoContent {
		t.Errorf("logout again = %d, without cookie = %d; want 204 both", again.Code, without.Code)
	}
	if me.Code != http.StatusUnauthorized {
		t.Errorf("me after logout = %d, want 401", me.Code)
	}
}

func TestDeleteMe(t *testing.T) {
	app := newAuthApp(t)
	cookie := app.signIn(t)

	deleted := app.do(http.MethodDelete, "/api/me", "", cookie)
	again := app.do(http.MethodDelete, "/api/me", "", cookie)
	anonymous := app.do(http.MethodDelete, "/api/me", "")
	login := app.do(http.MethodPost, "/api/auth/login", `{"email":"ana@example.com"}`)

	cleared := sessionCookieOf(deleted)
	if deleted.Code != http.StatusNoContent || cleared == nil || cleared.MaxAge >= 0 {
		t.Errorf("delete me = %d %s, cookie %v; want 204 with the cookie cleared", deleted.Code, deleted.Body.String(), cleared)
	}
	if again.Code != http.StatusUnauthorized || anonymous.Code != http.StatusUnauthorized {
		t.Errorf("delete me again = %d, without session = %d; want 401 both", again.Code, anonymous.Code)
	}
	if codes := app.codesSent(); login.Code != http.StatusNoContent || len(codes) != 1 {
		t.Errorf("login after deleting = %d with %d codes sent; want 204 and no new code (the account is gone)", login.Code, len(codes))
	}
}
