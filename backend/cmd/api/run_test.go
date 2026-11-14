package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

const runStopTimeout = 5 * time.Second

type runningApp struct {
	port   int
	dbPath string
	logs   *bytes.Buffer // read only after run returned
	cancel context.CancelFunc
	done   chan error
}

func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func startApp(t *testing.T) runningApp {
	t.Helper()
	app := runningApp{
		port:   freePort(t),
		dbPath: filepath.Join(t.TempDir(), "data", "app.db"),
		logs:   &bytes.Buffer{},
		done:   make(chan error, 1),
	}
	env := fakeEnv(map[string]string{
		"MAILCHECKER_PORT":    strconv.Itoa(app.port),
		"MAILCHECKER_DB_PATH": app.dbPath,
	})
	ctx, cancel := context.WithCancel(context.Background())
	app.cancel = cancel
	go func() { app.done <- run(ctx, env, app.logs) }()
	t.Cleanup(cancel)
	return app
}

func (app runningApp) healthStatus() (int, error) {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/health", app.port))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func (app runningApp) waitUntilHealthy(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(runStopTimeout)
	for time.Now().Before(deadline) {
		if status, err := app.healthStatus(); err == nil && status == http.StatusOK {
			return
		}
		select {
		case err := <-app.done:
			t.Fatalf("run() returned before serving: %v", err)
		case <-time.After(20 * time.Millisecond):
		}
	}
	t.Fatalf("GET /api/health never returned 200 within %s", runStopTimeout)
}

func (app runningApp) stop(t *testing.T) error {
	t.Helper()
	app.cancel()
	select {
	case err := <-app.done:
		return err
	case <-time.After(runStopTimeout):
		t.Fatalf("run() did not return within %s after cancel", runStopTimeout)
		return nil
	}
}

func TestRun_MigratesDBAndServesHealth(t *testing.T) {
	app := startApp(t)
	app.waitUntilHealthy(t)
	if err := app.stop(t); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	db, err := sqlite.Open(context.Background(), app.dbPath)
	if err != nil {
		t.Fatalf("reopen db: %v", err)
	}
	defer db.Close()
	var tables int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE name = 'goose_db_version'`).Scan(&tables); err != nil {
		t.Fatalf("query migrations table: %v", err)
	}
	if tables != 1 {
		t.Errorf("goose_db_version tables = %d, want 1 (db migrated)", tables)
	}
}

func TestRun_StopsCleanlyWhenContextCanceled(t *testing.T) {
	app := startApp(t)
	app.waitUntilHealthy(t)

	if err := app.stop(t); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}
	// SQLite removes the WAL file only when the last connection closes cleanly.
	if _, err := os.Stat(app.dbPath + "-wal"); !os.IsNotExist(err) {
		t.Errorf("%s-wal still exists (stat err = %v), want db closed", app.dbPath, err)
	}
}

func TestRun_FailsOnInvalidConfigWithoutOpeningDB(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app.db")
	env := fakeEnv(map[string]string{
		"MAILCHECKER_PORT":    "abc",
		"MAILCHECKER_DB_PATH": dbPath,
	})

	err := run(context.Background(), env, &bytes.Buffer{})

	if err == nil {
		t.Fatalf("run() error = nil, want config error")
	}
	if _, statErr := os.Stat(dbPath); !os.IsNotExist(statErr) {
		t.Errorf("db file exists (stat err = %v), want not opened", statErr)
	}
}

func TestRun_LogsPortAndDBPathWithoutSecrets(t *testing.T) {
	app := startApp(t)
	app.waitUntilHealthy(t)
	if err := app.stop(t); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	logs := app.logs.String()
	if !strings.Contains(logs, "port="+strconv.Itoa(app.port)) {
		t.Errorf("logs %q do not contain port=%d", logs, app.port)
	}
	if !strings.Contains(logs, app.dbPath) {
		t.Errorf("logs %q do not contain db path %q", logs, app.dbPath)
	}
	// The temp db path contains this test's name ("...WithoutSecrets"), so it is
	// removed before scanning for secret-like words.
	logsWithoutDBPath := strings.ToLower(strings.ReplaceAll(logs, app.dbPath, ""))
	for _, secret := range []string{"password", "MAILCHECKER_KEY", "secret"} {
		if strings.Contains(logsWithoutDBPath, strings.ToLower(secret)) {
			t.Errorf("logs contain %q, want no secrets", secret)
		}
	}
}

func TestRun_FailsWhenAMigrationFails(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app.db")
	db, err := sqlite.Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// A users table already there makes the migration that creates it fail.
	if _, err := db.Exec(`CREATE TABLE users (id INTEGER)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	db.Close()
	port := freePort(t)
	env := fakeEnv(map[string]string{"MAILCHECKER_PORT": strconv.Itoa(port), "MAILCHECKER_DB_PATH": dbPath})

	err = run(context.Background(), env, &bytes.Buffer{})

	if err == nil || !strings.Contains(err.Error(), "migrat") {
		t.Fatalf("run() error = %v, want the migration error", err)
	}
	if _, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/health", port)); err == nil {
		t.Errorf("the app answered after a failed migration")
	}
}
