package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	subsqlite "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/ports"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

const streamingCancelURL = "https://streaming.example/cancelar"

// fakeCatalog knows one service with a cancel link; any other merchant is unknown.
type fakeCatalog struct{}

func (fakeCatalog) Identify(merchant string) (domain.Service, bool) {
	if merchant == "STREAMING CO" {
		return domain.Service{Name: "Streaming Co", CancelURL: streamingCancelURL}, true
	}
	return domain.Service{}, false
}

func (fakeCatalog) Normalize(merchant string) string { return merchant }

func (fakeCatalog) LooksLikeSubscription(string) bool { return false }

var _ ports.ServiceCatalog = fakeCatalog{}

// testServer serves the subscriptions routes of one signed-in user over a real database.
type testServer struct {
	router http.Handler
	db     *sql.DB
	userID int64
}

func newTestServer(t *testing.T) testServer {
	t.Helper()
	db := sqlitetest.OpenMigrated(t)
	userID := sqlitetest.Users{DB: db}.Insert(t, "ana@example.com")
	responder := httpserver.Responder{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	router := httpserver.NewRouter(responder.Logger, httpserver.AlwaysHealthy)
	router.Group(func(protected chi.Router) {
		protected.Use(signedInAs(userID))
		handlersOver(db, responder).MountProtected(protected)
	})
	return testServer{router: router, db: db, userID: userID}
}

func signedInAs(userID int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(httpserver.WithUserID(r.Context(), userID)))
		})
	}
}

func handlersOver(db *sql.DB, responder httpserver.Responder) Handlers {
	now := func() time.Time { return asOf }
	lines := usecases.ChargeLines{Charges: subsqlite.NewChargeReader(db), Catalog: fakeCatalog{}, Now: now}
	decisions := subsqlite.NewDecisions(db, now)
	aliases := subsqlite.NewAliases(db)
	manuals := subsqlite.NewManuals(db)
	overrides := subsqlite.NewOverrides(db)
	keptRises := subsqlite.NewPriceAlertDecisions(db)
	hidden := subsqlite.NewHiddenLines(db)
	cancellations := subsqlite.NewCancellations(db)
	sources := usecases.SummarySources{
		Lines: lines, Decisions: decisions, Aliases: aliases, Rates: subsqlite.NewExchangeRates(db),
		KeptRises: keptRises, Manuals: manuals, Overrides: overrides, Hidden: hidden, Cancellations: cancellations,
	}
	return Handlers{
		Summary:             usecases.NewGetSummary(sources),
		Decide:              usecases.NewDecidePossibleSubscription(lines, decisions),
		UndoDecision:        usecases.NewUndoDecision(decisions),
		NameLine:            usecases.NewNameLine(lines, aliases),
		SaveOverride:        usecases.NewSaveLineOverride(lines, overrides),
		CreateManual:        usecases.NewCreateManualSubscription(manuals, sequentialIDs()),
		EditManual:          usecases.NewEditManualSubscription(manuals),
		RemoveManual:        usecases.NewRemoveManualSubscription(manuals),
		KeepRise:            usecases.NewKeepPriceRise(lines, keptRises),
		CancelGuide:         usecases.NewGetCancelGuide(lines),
		MarkForCancellation: usecases.NewMarkForCancellation(lines, cancellations),
		StopFollowing:       usecases.NewStopFollowingLine(lines, hidden),
		Responder:           responder,
	}
}

func sequentialIDs() func() string {
	next := 0
	return func() string {
		next++
		return fmt.Sprintf("0190a000-0000-7000-8000-%012d", next)
	}
}

// seedLines gives the user a known service whose last charge went up (active,
// with a rise) and an unknown merchant seen twice (possible). The possible
// merchant has the characters a key can carry in a URL.
func (s testServer) seedLines(t *testing.T) {
	t.Helper()
	for i, charge := range []struct {
		merchant string
		amount   shared.Money
		daysAgo  int
	}{
		{"STREAMING CO", 1000000, 70}, {"STREAMING CO", 1000000, 40}, {"STREAMING CO", 1050000, 10},
		{"MUSIC #1 APP", 399000, 40}, {"MUSIC #1 APP", 399000, 10},
	} {
		s.insertCharge(t, seededCharge{index: i, merchant: charge.merchant, amount: charge.amount, occurredAt: asOf.AddDate(0, 0, -charge.daysAgo)})
	}
}

// seededCharge is one charge in its own statement; the index keeps billing dates apart.
type seededCharge struct {
	index      int
	merchant   string
	amount     shared.Money
	occurredAt time.Time
}

func (s testServer) insertCharge(t *testing.T, charge seededCharge) {
	t.Helper()
	externalID := fmt.Sprintf("st:%d", charge.index)
	result, err := s.db.Exec(`INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at) VALUES (?, ?, 'bci', 'c.pdf', 'pdf', ?, 1)`,
		externalID, s.userID, charge.occurredAt.Unix()+int64(charge.index))
	if err != nil {
		t.Fatalf("insert statement: %v", err)
	}
	statementID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("statement id: %v", err)
	}
	_, err = s.db.Exec(`INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant, card, card_last4) VALUES (?, 'bci', 'statement', ?, ?, ?, ?, 'CLP', ?, 'credito', '1234')`,
		s.userID, externalID, statementID, charge.occurredAt.Unix(), charge.amount, charge.merchant)
	if err != nil {
		t.Fatalf("insert charge: %v", err)
	}
}

func (s testServer) send(t *testing.T, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, httptest.NewRequest(method, target, strings.NewReader(body)))
	return recorder
}

// summaryView is the part of GET /api/summary the action tests read back.
type summaryView struct {
	Subscriptions []lineView `json:"subscriptions"`
	Possible      []lineView `json:"possible"`
	Rises         []struct {
		LineKey string `json:"lineKey"`
	} `json:"rises"`
}

type lineView struct {
	LineKey                 string `json:"lineKey"`
	Alias                   string `json:"alias"`
	MarkedForCancellationAt *int64 `json:"markedForCancellationAt"`
	Amount                  struct {
		Amount int64 `json:"amount"`
	} `json:"amount"`
}

func (s testServer) summary(t *testing.T) summaryView {
	t.Helper()
	recorder := s.send(t, http.MethodGet, "/api/summary", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /api/summary = %d %s", recorder.Code, recorder.Body.String())
	}
	var view summaryView
	if err := json.Unmarshal(recorder.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	return view
}

// seededKeys are the line keys of seedLines, read from the summary itself.
type seededKeys struct {
	active   string
	possible string
}

func (s testServer) seededKeys(t *testing.T) seededKeys {
	t.Helper()
	view := s.summary(t)
	if len(view.Subscriptions) != 1 || len(view.Possible) != 1 {
		t.Fatalf("seeded summary = %+v, want one active and one possible line", view)
	}
	return seededKeys{active: view.Subscriptions[0].LineKey, possible: view.Possible[0].LineKey}
}

func (s testServer) find(t *testing.T, key string) (lineView, bool) {
	t.Helper()
	view := s.summary(t)
	for _, line := range append(view.Subscriptions, view.Possible...) {
		if line.LineKey == key {
			return line, true
		}
	}
	return lineView{}, false
}

func assertResponse(t *testing.T, recorder *httptest.ResponseRecorder, status int, bodyHas string) {
	t.Helper()
	if recorder.Code != status || !strings.Contains(recorder.Body.String(), bodyHas) {
		t.Errorf("response = %d %s, want %d containing %s", recorder.Code, recorder.Body.String(), status, bodyHas)
	}
}

func jsonBody(t *testing.T, value any) string {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}
	return string(body)
}

// insertManualOf stores a manual of another user, to check it is out of reach.
func insertManualOf(t *testing.T, db *sql.DB, email string) string {
	t.Helper()
	otherID := sqlitetest.Users{DB: db}.Insert(t, email)
	manual := domain.ManualSubscription{PublicID: "0190a000-0000-7000-8000-999999999999", Name: "Ajena", Amount: 100000, Currency: shared.CurrencyCLP, BillingDay: 5}
	if err := subsqlite.NewManuals(db).Create(context.Background(), otherID, manual); err != nil {
		t.Fatalf("insert manual of %s: %v", email, err)
	}
	return manual.PublicID
}
