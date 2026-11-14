package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite"
)

// Smoke tests run the real app (run, real database, console codes) and go
// through it only over HTTP, like the front will.

// syncBuffer is the app's stdout, read while the app still writes to it.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

type smokeApp struct {
	base   string
	dbPath string
	out    *syncBuffer
}

func startSmokeApp(t *testing.T) smokeApp {
	t.Helper()
	port := freePort(t)
	app := smokeApp{base: fmt.Sprintf("http://127.0.0.1:%d", port), dbPath: filepath.Join(t.TempDir(), "app.db"), out: &syncBuffer{}}
	env := fakeEnv(map[string]string{"MAILCHECKER_PORT": strconv.Itoa(port), "MAILCHECKER_DB_PATH": app.dbPath})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- run(ctx, env, app.out) }()
	t.Cleanup(func() {
		// A client may keep a spare connection it never used; the server waits
		// for it on shutdown, so the test's clients let go of theirs first.
		http.DefaultTransport.(*http.Transport).CloseIdleConnections()
		cancel()
		if err := <-done; err != nil {
			t.Errorf("run() error = %v", err)
		}
	})
	app.waitUntilHealthy(t)
	return app
}

func (a smokeApp) waitUntilHealthy(t *testing.T) {
	t.Helper()
	for deadline := time.Now().Add(runStopTimeout); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if res, err := http.Get(a.base + "/api/health"); err == nil {
			res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return
			}
		}
	}
	t.Fatalf("the app never got healthy")
}

// lastCode is the last code the console printed for that email.
func (a smokeApp) lastCode(t *testing.T, email string) string {
	t.Helper()
	matches := regexp.MustCompile(`login code for `+regexp.QuoteMeta(email)+`: (\d{6})`).FindAllStringSubmatch(a.out.String(), -1)
	if len(matches) == 0 {
		t.Fatalf("no code printed for %s", email)
	}
	return matches[len(matches)-1][1]
}

// smokeUser is a browser of one person: their cookie jar.
type smokeUser struct {
	app    smokeApp
	client *http.Client
	email  string
}

type smokeReply struct {
	status int
	body   string
}

func (a smokeApp) newUser(t *testing.T, email string) smokeUser {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookie jar: %v", err)
	}
	return smokeUser{app: a, client: &http.Client{Jar: jar}, email: email}
}

func (u smokeUser) do(t *testing.T, method, path, body string) smokeReply {
	t.Helper()
	request, err := http.NewRequest(method, u.app.base+path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	res, err := u.client.Do(request)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read %s %s: %v", method, path, err)
	}
	return smokeReply{status: res.StatusCode, body: string(read)}
}

func (u smokeUser) expect(t *testing.T, reply smokeReply, status int) smokeReply {
	t.Helper()
	if reply.status != status {
		t.Fatalf("%s got %d %s, want %d", u.email, reply.status, reply.body, status)
	}
	return reply
}

func (u smokeUser) register(t *testing.T, name string) {
	t.Helper()
	u.expect(t, u.do(t, http.MethodPost, "/api/auth/register", fmt.Sprintf(`{"name":%q,"email":%q,"acceptedTerms":true}`, name, u.email)), http.StatusNoContent)
}

func (u smokeUser) verify(t *testing.T, code string) smokeReply {
	t.Helper()
	return u.do(t, http.MethodPost, "/api/auth/code/verify", fmt.Sprintf(`{"email":%q,"code":%q}`, u.email, code))
}

// signUp registers and opens a session with the printed code.
func (a smokeApp) signUp(t *testing.T, name, email string) smokeUser {
	t.Helper()
	user := a.newUser(t, email)
	user.register(t, name)
	user.expect(t, user.verify(t, a.lastCode(t, email)), http.StatusNoContent)
	return user
}

func TestSmoke_Access(t *testing.T) {
	app := startSmokeApp(t)
	ana := app.signUp(t, "Ana Prueba", "ana@example.com")

	me := ana.expect(t, ana.do(t, http.MethodGet, "/api/me", ""), http.StatusOK)
	if !strings.Contains(me.body, `"name":"Ana Prueba"`) {
		t.Errorf("me = %s, want her name", me.body)
	}
	ana.expect(t, ana.do(t, http.MethodPost, "/api/auth/logout", ""), http.StatusNoContent)
	ana.expect(t, ana.do(t, http.MethodGet, "/api/me", ""), http.StatusUnauthorized)
}

func TestSmoke_FiveWrongCodesBurnTheRightOne(t *testing.T) {
	app := startSmokeApp(t)
	beto := app.newUser(t, "beto@example.com")
	beto.register(t, "Beto Prueba")
	right := app.lastCode(t, beto.email)
	wrong := "000000"
	if right == wrong {
		wrong = "111111"
	}

	for range 5 {
		beto.verify(t, wrong)
	}

	if got := beto.verify(t, right); got.status == http.StatusNoContent {
		t.Errorf("right code after five wrong ones = %d, want it rejected", got.status)
	}
	beto.expect(t, beto.do(t, http.MethodGet, "/api/me", ""), http.StatusUnauthorized)
}

// seedPossible gives the user a merchant charged twice a month apart: a
// possible subscription. Statements go straight to the app's database file.
func (a smokeApp) seedPossible(t *testing.T, email, merchant string) {
	t.Helper()
	db := a.openDB(t)
	defer db.Close()
	now := time.Now().UTC()
	for i, daysAgo := range []int{40, 10} {
		at := now.AddDate(0, 0, -daysAgo).Unix()
		publicID := fmt.Sprintf("%s-%s-%d", email, merchant, i)
		result, err := db.Exec(`INSERT INTO statements (public_id, user_id, bank, filename, kind, billing_date, imported_at)
			VALUES (?, (SELECT id FROM users WHERE email = ?), 'bci', 'c.pdf', 'pdf', ?, ?)`, publicID, email, at, at)
		if err != nil {
			t.Fatalf("seed statement: %v", err)
		}
		statementID, err := result.LastInsertId()
		if err != nil {
			t.Fatalf("statement id: %v", err)
		}
		if _, err := db.Exec(`INSERT INTO charges (user_id, bank, source, external_id, statement_id, occurred_at, amount, currency, merchant)
			VALUES ((SELECT id FROM users WHERE email = ?), 'bci', 'statement', ?, ?, ?, 399000, 'CLP', ?)`, email, publicID, statementID, at, merchant); err != nil {
			t.Fatalf("seed charge: %v", err)
		}
	}
}

func (a smokeApp) openDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sqlite.Open(context.Background(), a.dbPath)
	if err != nil {
		t.Fatalf("open app database: %v", err)
	}
	return db
}

type smokeSummary struct {
	Possible []struct {
		LineKey string `json:"lineKey"`
		Name    string `json:"name"`
	} `json:"possible"`
}

func (u smokeUser) possibleKey(t *testing.T) string {
	t.Helper()
	var summary smokeSummary
	if err := json.Unmarshal([]byte(u.expect(t, u.do(t, http.MethodGet, "/api/summary", ""), http.StatusOK).body), &summary); err != nil || len(summary.Possible) != 1 {
		t.Fatalf("summary of %s = %+v (%v), want one possible", u.email, summary, err)
	}
	return summary.Possible[0].LineKey
}

func (u smokeUser) createManual(t *testing.T) string {
	t.Helper()
	var created struct {
		ID string `json:"id"`
	}
	reply := u.expect(t, u.do(t, http.MethodPost, "/api/subscriptions/manual", `{"name":"Gimnasio","amount":2500000,"currency":"CLP","billingDay":15}`), http.StatusCreated)
	if err := json.Unmarshal([]byte(reply.body), &created); err != nil {
		t.Fatalf("created manual = %s: %v", reply.body, err)
	}
	return created.ID
}

func TestSmoke_UsersDoNotSeeEachOther(t *testing.T) {
	app := startSmokeApp(t)
	ana := app.signUp(t, "Ana Prueba", "ana@example.com")
	beto := app.signUp(t, "Beto Prueba", "beto@example.com")
	app.seedPossible(t, ana.email, "MUSIC APP")
	app.seedPossible(t, beto.email, "VIDEO APP")

	anaKey := ana.possibleKey(t)
	betoSummary := beto.expect(t, beto.do(t, http.MethodGet, "/api/summary", ""), http.StatusOK).body
	if strings.Contains(betoSummary, "MUSIC APP") || !strings.Contains(betoSummary, "VIDEO APP") {
		t.Errorf("summary of beto = %s, want only his merchant", betoSummary)
	}
	beto.expect(t, beto.do(t, http.MethodPost, "/api/subscriptions/decisions", fmt.Sprintf(`{"groupKey":%q,"decision":"dismiss"}`, anaKey)), http.StatusNotFound)
	anaManual := ana.createManual(t)
	beto.expect(t, beto.do(t, http.MethodPut, "/api/subscriptions/manual/"+anaManual, `{"name":"Mía","amount":1,"currency":"CLP","billingDay":1}`), http.StatusNotFound)
	beto.expect(t, beto.do(t, http.MethodDelete, "/api/subscriptions/manual/"+anaManual, ""), http.StatusNotFound)
	if key := ana.possibleKey(t); key != anaKey {
		t.Errorf("possible of ana = %q, want %q untouched", key, anaKey)
	}
}

// rowsOfUser counts the user's rows in every table that holds their data.
func (a smokeApp) rowsOfUser(t *testing.T, email string) int {
	t.Helper()
	db := a.openDB(t)
	defer db.Close()
	rows, err := db.Query(`SELECT m.name FROM sqlite_master m JOIN pragma_table_info(m.name) c ON c.name = 'user_id' WHERE m.type = 'table'`)
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan: %v", err)
		}
		tables = append(tables, name)
	}
	rows.Close()
	total := a.countOf(t, db, "SELECT count(*) FROM users WHERE email = ?", email) + a.countOf(t, db, "SELECT count(*) FROM code_sends WHERE email = ?", email)
	for _, table := range tables {
		total += a.countOf(t, db, "SELECT count(*) FROM "+table+" WHERE user_id = (SELECT id FROM users WHERE email = ?)", email)
	}
	return total
}

func (a smokeApp) countOf(t *testing.T, db *sql.DB, query, email string) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, email).Scan(&n); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return n
}

// withChoices gives the user a possible line, a decision, an alias and a manual.
func (a smokeApp) withChoices(t *testing.T, user smokeUser) {
	t.Helper()
	a.seedPossible(t, user.email, "MUSIC APP")
	key := user.possibleKey(t)
	user.expect(t, user.do(t, http.MethodPut, "/api/subscriptions/aliases", fmt.Sprintf(`{"lineKey":%q,"alias":"Casa"}`, key)), http.StatusNoContent)
	user.expect(t, user.do(t, http.MethodPost, "/api/subscriptions/decisions", fmt.Sprintf(`{"groupKey":%q,"decision":"confirm","cadence":"monthly"}`, key)), http.StatusNoContent)
	user.createManual(t)
}

func TestSmoke_DeleteAll(t *testing.T) {
	app := startSmokeApp(t)
	ana := app.signUp(t, "Ana Prueba", "ana@example.com")
	beto := app.signUp(t, "Beto Prueba", "beto@example.com")
	app.withChoices(t, ana)
	app.withChoices(t, beto)
	betoBefore := app.rowsOfUser(t, beto.email)

	ana.expect(t, ana.do(t, http.MethodDelete, "/api/me", ""), http.StatusNoContent)

	if left := app.rowsOfUser(t, ana.email); left != 0 {
		t.Errorf("rows of ana after deleting = %d, want 0", left)
	}
	ana.expect(t, ana.do(t, http.MethodGet, "/api/summary", ""), http.StatusUnauthorized)
	if after := app.rowsOfUser(t, beto.email); after != betoBefore {
		t.Errorf("rows of beto = %d, want %d untouched", after, betoBefore)
	}
	beto.expect(t, beto.do(t, http.MethodGet, "/api/me", ""), http.StatusOK)
}
