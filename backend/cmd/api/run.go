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

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/migrations"
)

const (
	listenHost        = "127.0.0.1" // personal, local app: never exposed to the network
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 5 * time.Minute // an email sync can take a while
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 5 * time.Second
)

// run wires config → logger → database → migrations → router → server and
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
	listener, err := listen(config.Port)
	if err != nil {
		return err
	}
	logger.Info("server started", "port", config.Port, "db_path", config.DBPath)
	return serveUntilCanceled(ctx, newServer(httpserver.NewRouter(logger)), listener, logger)
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

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

func serveUntilCanceled(ctx context.Context, server *http.Server, listener net.Listener, logger *slog.Logger) error {
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(listener) }()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}
	if err := shutdown(server); err != nil {
		return err
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}
	logger.Info("server stopped")
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
