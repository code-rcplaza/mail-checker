package main

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/export"
	authhttp "github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/http"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/mailer"
	authsqlite "github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite"
	authports "github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	authusecases "github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/clock"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/bci"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/extract"
	statementshttp "github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/http"
	statementssqlite "github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/sqlite"
	statementsports "github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
	statementsusecases "github.com/code-rcplaza/mail-checker/backend/internal/statements/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/catalog"
	subscriptionshttp "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/http"
	subscriptionssqlite "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/adapters/sqlite"
	subscriptionsusecases "github.com/code-rcplaza/mail-checker/backend/internal/subscriptions/usecases"
)

const (
	// userLocation is the calendar of the users: statement dates and months are read in it.
	userLocation  = "America/Santiago"
	resendTimeout = 10 * time.Second
)

// appDependencies are what the composition needs from run.
type appDependencies struct {
	db     *sql.DB
	config Config
	logger *slog.Logger
	stdout io.Writer // where the console sender prints codes in development
}

// newPublicID is a UUID v7. NewV7 only fails if the system has no randomness,
// and then nothing in the app can work.
func newPublicID() string { return uuid.Must(uuid.NewV7()).String() }

// buildRouter creates every repository and use case once and mounts the
// routes: health and auth are public, the rest goes behind the session.
func buildRouter(deps appDependencies) (http.Handler, error) {
	location, err := time.LoadLocation(userLocation)
	if err != nil {
		return nil, fmt.Errorf("load location %s: %w", userLocation, err)
	}
	sender, err := codeSenderFor(deps.config, deps.stdout)
	if err != nil {
		return nil, err
	}
	responder := httpserver.Responder{Logger: deps.logger}
	sessions := authsqlite.NewSessions(deps.db)
	auth := authHandlers(deps, authWiring{sender: sender, sessions: sessions, responder: responder})
	statements := statementsHandlers(deps.db, statementsWiring{location: location, responder: responder})
	subscriptions := subscriptionsHandlers(deps.db, responder)

	router := httpserver.NewRouter(deps.logger, deps.db.PingContext)
	auth.MountPublic(router)
	router.Group(func(protected chi.Router) {
		protected.Use(authhttp.RequireSession{Validate: authusecases.NewValidateSession(sessions, clock.Now), Responder: responder}.Handler)
		auth.MountProtected(protected)
		statements.MountProtected(protected)
		subscriptions.MountProtected(protected)
	})
	return router, nil
}

// codeSenderFor mails codes with Resend whenever there is a key; only
// development without one prints them, and production never does.
func codeSenderFor(config Config, stdout io.Writer) (authports.CodeSender, error) {
	if !config.Mailer.UsesConsole() {
		return mailer.NewResendSender(mailer.ResendSettings{
			BaseURL: mailer.ResendBaseURL, APIKey: string(config.Mailer.ResendAPIKey), From: config.Mailer.MailFrom, Timeout: resendTimeout,
		}), nil
	}
	if config.IsProduction {
		return nil, fmt.Errorf("%s: production cannot print login codes", envResendKey)
	}
	return mailer.ConsoleSender{Writer: stdout}, nil
}

// bankParsers is the registry of banks: built once here and only read after.
func bankParsers(location *time.Location) []statementsports.BankStatementParser {
	return []statementsports.BankStatementParser{bci.NewStatementParser(location)}
}

type authWiring struct {
	sender    authports.CodeSender
	sessions  authsqlite.Sessions
	responder httpserver.Responder
}

func authHandlers(deps appDependencies, wiring authWiring) authhttp.Handlers {
	users := authsqlite.NewUsers(deps.db, clock.Now)
	codes := authsqlite.NewLoginCodes(deps.db)
	delivery := authusecases.NewCodeDelivery(authusecases.CodeDeliveryDependencies{
		Codes: codes, Sends: authsqlite.NewCodeSends(deps.db), Sender: wiring.sender, Random: cryptoRandom, Now: clock.Now,
	})
	verify := authusecases.NewVerifyCode(authusecases.VerifyCodeDependencies{Users: users, Codes: codes, Now: clock.Now})
	return authhttp.Handlers{
		Register:     authusecases.NewRegisterUser(authusecases.RegisterUserDependencies{Users: users, Delivery: delivery, NewPublicID: newPublicID, Now: clock.Now}),
		Login:        authusecases.NewLogin(users, delivery),
		OpenSession:  authusecases.NewOpenSession(authusecases.OpenSessionDependencies{Verify: verify, Sessions: wiring.sessions, Random: cryptoRandom, Now: clock.Now}),
		CloseSession: authusecases.NewCloseSession(wiring.sessions),
		GetMe:        authusecases.NewGetMe(users),
		DeleteUser:   authusecases.NewDeleteUser(users),
		Export:       userCopies(deps.db, users),
		Cookies:      authhttp.SessionCookies{IsProduction: deps.config.IsProduction},
		Responder:    wiring.responder,
	}
}

type statementsWiring struct {
	location  *time.Location
	responder httpserver.Responder
}

func statementsHandlers(db *sql.DB, wiring statementsWiring) statementshttp.Handlers {
	parsers := bankParsers(wiring.location)
	return statementshttp.Handlers{
		Import: statementsusecases.NewImportStatement(statementsusecases.ImportStatementDependencies{
			Extractor: extract.PDFExtractor{}, Parsers: parsers, Statements: statementssqlite.NewStatements(db),
			Rates: statementssqlite.NewExchangeRates(db), NewPublicID: newPublicID, Now: clock.Now, Location: wiring.location,
		}),
		Banks:     statementsusecases.NewListBanks(parsers),
		Limits:    statementshttp.DefaultUploadLimits,
		Responder: wiring.responder,
	}
}

// subscriptionRepos are the subscriptions repositories, shared by the routes and the copy.
type subscriptionRepos struct {
	lines         subscriptionsusecases.ChargeLines
	decisions     subscriptionssqlite.Decisions
	aliases       subscriptionssqlite.Aliases
	manuals       subscriptionssqlite.Manuals
	overrides     subscriptionssqlite.Overrides
	keptRises     subscriptionssqlite.PriceAlertDecisions
	hidden        subscriptionssqlite.HiddenLines
	cancellations subscriptionssqlite.Cancellations
}

func subscriptionReposOver(db *sql.DB) subscriptionRepos {
	return subscriptionRepos{
		lines:     subscriptionsusecases.ChargeLines{Charges: subscriptionssqlite.NewChargeReader(db), Catalog: catalog.Catalog{}, Now: clock.Now},
		decisions: subscriptionssqlite.NewDecisions(db, clock.Now), aliases: subscriptionssqlite.NewAliases(db),
		manuals: subscriptionssqlite.NewManuals(db), overrides: subscriptionssqlite.NewOverrides(db),
		keptRises: subscriptionssqlite.NewPriceAlertDecisions(db), hidden: subscriptionssqlite.NewHiddenLines(db),
		cancellations: subscriptionssqlite.NewCancellations(db),
	}
}

func summaryOver(db *sql.DB, repos subscriptionRepos) subscriptionsusecases.GetSummary {
	return subscriptionsusecases.NewGetSummary(subscriptionsusecases.SummarySources{
		Lines: repos.lines, Decisions: repos.decisions, Aliases: repos.aliases, Rates: subscriptionssqlite.NewExchangeRates(db),
		KeptRises: repos.keptRises, Manuals: repos.manuals, Overrides: repos.overrides, Hidden: repos.hidden, Cancellations: repos.cancellations,
	})
}

func subscriptionsHandlers(db *sql.DB, responder httpserver.Responder) subscriptionshttp.Handlers {
	repos := subscriptionReposOver(db)
	return subscriptionshttp.Handlers{
		Summary:             summaryOver(db, repos),
		Decide:              subscriptionsusecases.NewDecidePossibleSubscription(repos.lines, repos.decisions),
		UndoDecision:        subscriptionsusecases.NewUndoDecision(repos.decisions),
		NameLine:            subscriptionsusecases.NewNameLine(repos.lines, repos.aliases),
		SaveOverride:        subscriptionsusecases.NewSaveLineOverride(repos.lines, repos.overrides),
		CreateManual:        subscriptionsusecases.NewCreateManualSubscription(repos.manuals, newPublicID),
		EditManual:          subscriptionsusecases.NewEditManualSubscription(repos.manuals),
		RemoveManual:        subscriptionsusecases.NewRemoveManualSubscription(repos.manuals),
		KeepRise:            subscriptionsusecases.NewKeepPriceRise(repos.lines, repos.keptRises),
		CancelGuide:         subscriptionsusecases.NewGetCancelGuide(repos.lines),
		MarkForCancellation: subscriptionsusecases.NewMarkForCancellation(repos.lines, repos.cancellations),
		StopFollowing:       subscriptionsusecases.NewStopFollowingLine(repos.lines, repos.hidden),
		Responder:           responder,
	}
}

func userCopies(db *sql.DB, users authsqlite.Users) export.UserCopies {
	repos := subscriptionReposOver(db)
	return export.NewUserCopies(export.CopySources{
		Me:         authusecases.NewGetMe(users),
		Statements: statementsusecases.NewListStatements(statementssqlite.NewStatements(db)),
		Saved: subscriptionsusecases.NewListSavedChoices(subscriptionsusecases.SavedChoicesSources{
			Charges: repos.lines.Charges, Decisions: repos.decisions, Aliases: repos.aliases, Manuals: repos.manuals, Overrides: repos.overrides,
		}),
		Summary: summaryOver(db, repos),
	})
}
