package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	authsqlite "github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/sqlite"
	authusecases "github.com/code-rcplaza/mail-checker/backend/internal/auth/usecases"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/clock"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/migrations"
)

const (
	listenHost      = "127.0.0.1" // personal, local app: never exposed to the network
	shutdownTimeout = 5 * time.Second
)

// run wires config → logger → database → migrations → router → cleanup → server and
// blocks until ctx is canceled, then shuts down in order.
func run(ctx context.Context, getenv func(string) string, stdout io.Writer) error {
	config, err := LoadConfig(getenv)
	if err != nil {
		return err
	}
	logger := newLogger(stdout)
	db, err := sqlite.Open(ctx, config.DBPath)
	if err != nil {
		return err
	}
	defer closeDB(db, logger)
	if err := sqlite.Migrate(ctx, db, migrations.FS); err != nil {
		return err
	}
	router, err := buildRouter(appDependencies{db: db, config: config, logger: logger, stdout: stdout})
	if err != nil {
		return err
	}
	listener, err := listen(config.Port)
	if err != nil {
		return err
	}
	go cleanupEveryHour(ctx, appDependencies{db: db, logger: logger})
	logger.Info("server started", "port", config.Port, "db_path", config.DBPath, "production", config.IsProduction, "mailer", mailerName(config.Mailer))
	return runningServer{server: newServer(router, defaultTimeouts), listener: listener, logger: logger}.serveUntilCanceled(ctx)
}

func cleanupEveryHour(ctx context.Context, deps appDependencies) {
	cleanup := authusecases.NewCleanExpiredAccess(authsqlite.NewExpiredAccess(deps.db), clock.Now)
	cleanupLoop{clean: cleanup.Execute, every: cleanupInterval, logger: deps.logger}.run(ctx)
}

// mailerName says how codes leave, for the start log; never the key.
func mailerName(mailer MailerConfig) string {
	if mailer.UsesConsole() {
		return "console"
	}
	return "resend"
}

func newLogger(stdout io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(stdout, nil))
}

func closeDB(db *sql.DB, logger *slog.Logger) {
	if err := db.Close(); err != nil {
		logger.Error("close database failed", "err", err)
	}
}

// listen binds before logging "started", so a busy port fails fast.
func listen(port int) (net.Listener, error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(listenHost, strconv.Itoa(port)))
	if err != nil {
		return nil, fmt.Errorf("listen on port %d: %w", port, err)
	}
	return listener, nil
}

// serverTimeouts bound every phase of a connection, so a slow or stuck client
// cannot hold the server.
type serverTimeouts struct {
	readHeader time.Duration
	read       time.Duration
	write      time.Duration
	idle       time.Duration
}

var defaultTimeouts = serverTimeouts{
	readHeader: 5 * time.Second,
	read:       30 * time.Second,
	write:      5 * time.Minute, // an email sync can take a while
	idle:       60 * time.Second,
}

func newServer(handler http.Handler, timeouts serverTimeouts) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: timeouts.readHeader,
		ReadTimeout:       timeouts.read,
		WriteTimeout:      timeouts.write,
		IdleTimeout:       timeouts.idle,
	}
}

// runningServer is the HTTP server, the port it listens on and its log.
type runningServer struct {
	server   *http.Server
	listener net.Listener
	logger   *slog.Logger
}

func (s runningServer) serveUntilCanceled(ctx context.Context) error {
	serveErr := make(chan error, 1)
	go func() { serveErr <- s.server.Serve(s.listener) }()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}
	if err := shutdown(s.server); err != nil {
		return err
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}
	s.logger.Info("server stopped")
	return nil
}

func shutdown(server *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
