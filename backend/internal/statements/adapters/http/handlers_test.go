package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
	statementsqlite "github.com/code-rcplaza/mail-checker/backend/internal/statements/adapters/sqlite"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/ports"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/usecases"
)

// textExtractor reads the upload as plain text; "NOT PDF" is not a statement.
type textExtractor struct{ delay time.Duration }

func (e textExtractor) Extract(_ string, data []byte) (string, domain.FileKind, error) {
	time.Sleep(e.delay)
	if strings.HasPrefix(string(data), "NOT PDF") {
		return "", "", shared.Invalid(domain.KeyFileKindUnknown, nil)
	}
	return string(data), domain.FileKindPDF, nil
}

// lineParser reads "BANCO BCI" files: one "MERCHANT;pesos" charge per line, "???" is unreadable.
type lineParser struct{}

func (lineParser) Bank() shared.Bank       { return "bci" }
func (lineParser) Detect(text string) bool { return strings.Contains(text, "BANCO BCI") }
func (lineParser) BillingDate(text string) (time.Time, bool) {
	return time.Time{}, false
}
func (lineParser) TransferPair(string) (domain.TransferPair, bool) {
	return domain.TransferPair{}, false
}
func (lineParser) Parse(text string) ([]domain.StatementLine, []int) {
	var lines []domain.StatementLine
	var bad []int
	for i, row := range strings.Split(text, "\n") {
		merchant, pesos, isCharge := strings.Cut(row, ";")
		if row == "???" {
			bad = append(bad, i+1)
		}
		amount, err := strconv.Atoi(pesos)
		if !isCharge || err != nil {
			continue
		}
		lines = append(lines, domain.StatementLine{Number: i + 1, Charge: shared.ParsedCharge{Amount: shared.Money(amount), Currency: shared.CurrencyCLP, Merchant: merchant, Card: "credito"}})
	}
	return lines, bad
}

var _ ports.BankStatementParser = lineParser{}

var billing = time.Date(2026, 10, 25, 3, 0, 0, 0, time.UTC)

type statementsApp struct {
	router http.Handler
	userID int64
}

func newStatementsApp(t *testing.T, limits UploadLimits, extractor ports.FileExtractor) statementsApp {
	t.Helper()
	db := sqlitetest.OpenMigrated(t)
	userID := sqlitetest.Users{DB: db}.Insert(t, "ana@example.com")
	parsers := []ports.BankStatementParser{lineParser{}}
	ids := 0
	handlers := Handlers{
		Import: usecases.NewImportStatement(usecases.ImportStatementDependencies{
			Extractor: extractor, Parsers: parsers, Statements: statementsqlite.NewStatements(db), Rates: statementsqlite.NewExchangeRates(db),
			NewPublicID: func() string { ids++; return fmt.Sprintf("statement-%d", ids) }, Now: func() time.Time { return billing }, Location: time.UTC,
		}),
		Banks:     usecases.NewListBanks(parsers),
		Limits:    limits,
		Responder: httpserver.Responder{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))},
	}
	router := httpserver.NewRouter(handlers.Responder.Logger, httpserver.AlwaysHealthy)
	router.Group(func(protected chi.Router) {
		protected.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.WithContext(httpserver.WithUserID(r.Context(), userID)))
			})
		})
		handlers.MountProtected(protected)
	})
	return statementsApp{router: router, userID: userID}
}

var testLimits = UploadLimits{MaxBytes: 1 << 20, Timeout: 5 * time.Second}

// upload is one multipart form; files maps a field to the contents of each file.
type upload struct {
	files  [][]byte
	fields map[string]string
}

func (u upload) request(t *testing.T) *http.Request {
	t.Helper()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for _, content := range u.files {
		part, err := form.CreateFormFile("file", "cartola.pdf")
		if err != nil {
			t.Fatalf("form file: %v", err)
		}
		part.Write(content)
	}
	for name, value := range u.fields {
		form.WriteField(name, value)
	}
	form.Close()
	request := httptest.NewRequest(http.MethodPost, "/api/statements", &body)
	request.Header.Set("Content-Type", form.FormDataContentType())
	return request
}

func (a statementsApp) send(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.router.ServeHTTP(recorder, request)
	return recorder
}

const goodStatement = "BANCO BCI\nNETFLIX;599000\n???\nSPOTIFY;399000"

func withDate(fields map[string]string) map[string]string {
	if fields == nil {
		fields = map[string]string{}
	}
	fields["billingDate"] = strconv.FormatInt(billing.Unix(), 10)
	return fields
}

type replyBody struct {
	ID       string   `json:"id"`
	Filename string   `json:"filename"`
	Months   []string `json:"months"`
	Charges  int      `json:"charges"`
	BadLines int      `json:"badLines"`
	Error    struct {
		Key    string            `json:"key"`
		Fields map[string]string `json:"fields"`
	} `json:"error"`
}

func replyOf(t *testing.T, recorder *httptest.ResponseRecorder) replyBody {
	t.Helper()
	var reply replyBody
	if err := json.Unmarshal(recorder.Body.Bytes(), &reply); err != nil {
		t.Fatalf("body %q: %v", recorder.Body.String(), err)
	}
	return reply
}

func TestUpload_SavesTheStatementOfTheUser(t *testing.T) {
	app := newStatementsApp(t, testLimits, textExtractor{})

	got := app.send(upload{files: [][]byte{[]byte(goodStatement)}, fields: withDate(nil)}.request(t))

	reply := replyOf(t, got)
	if got.Code != http.StatusCreated || reply.ID != "statement-1" || reply.Filename != "cartola.pdf" || reply.Charges != 2 || reply.BadLines != 1 || fmt.Sprint(reply.Months) != "[2026-10]" {
		t.Errorf("upload = %d %s, want 201 with id, filename, months, 2 charges and 1 bad line", got.Code, got.Body.String())
	}
	if strings.Contains(got.Body.String(), "???") {
		t.Errorf("the response carries the text of a bad line: %s", got.Body.String())
	}
}

func TestUpload_Errors(t *testing.T) {
	cases := []struct {
		name    string
		upload  upload
		status  int
		key     string
		inField string
	}{
		{"a bank that is not detected", upload{files: [][]byte{[]byte("OTRO BANCO\nX;1")}, fields: withDate(nil)}, http.StatusBadRequest, "statements.bank_not_detected", ""},
		{"a file without charges", upload{files: [][]byte{[]byte("BANCO BCI\n???")}, fields: withDate(nil)}, http.StatusBadRequest, "statements.no_charges", ""},
		{"a file that is not a PDF", upload{files: [][]byte{[]byte("NOT PDF")}, fields: withDate(nil)}, http.StatusBadRequest, "statements.file_kind_unknown", ""},
		{"an unknown bank", upload{files: [][]byte{[]byte(goodStatement)}, fields: withDate(map[string]string{"bank": "otro"})}, http.StatusBadRequest, "statements.bank_unknown", ""},
		{"a billing date that is not a number", upload{files: [][]byte{[]byte(goodStatement)}, fields: map[string]string{"billingDate": "ayer"}}, http.StatusBadRequest, "common.validation_failed", "billingDate"},
		{"no file", upload{fields: withDate(nil)}, http.StatusBadRequest, "common.validation_failed", "file"},
		{"two files", upload{files: [][]byte{[]byte(goodStatement), []byte(goodStatement)}, fields: withDate(nil)}, http.StatusBadRequest, "statements.file_too_many", ""},
		{"a file over the limit", upload{files: [][]byte{bytes.Repeat([]byte("a"), (1<<20)+1)}, fields: withDate(nil)}, http.StatusBadRequest, "statements.file_too_large", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := newStatementsApp(t, testLimits, textExtractor{}).send(tc.upload.request(t))

			reply := replyOf(t, got)
			if got.Code != tc.status || reply.Error.Key != tc.key {
				t.Errorf("upload = %d %s, want %d %s", got.Code, got.Body.String(), tc.status, tc.key)
			}
			if tc.inField != "" && reply.Error.Fields[tc.inField] == "" {
				t.Errorf("fields = %v, want %s", reply.Error.Fields, tc.inField)
			}
		})
	}
}

func TestUpload_AFileUnderTheLimitPasses(t *testing.T) {
	content := goodStatement + "\n" + strings.Repeat("x", 900_000)

	got := newStatementsApp(t, testLimits, textExtractor{}).send(upload{files: [][]byte{[]byte(content)}, fields: withDate(nil)}.request(t))

	if got.Code != http.StatusCreated {
		t.Errorf("upload of ~1 MB = %d %s, want 201", got.Code, got.Body.String())
	}
}

func TestUpload_TheSameStatementTwiceIsAConflict(t *testing.T) {
	app := newStatementsApp(t, testLimits, textExtractor{})
	first := upload{files: [][]byte{[]byte(goodStatement)}, fields: withDate(nil)}
	app.send(first.request(t))

	got := app.send(first.request(t))

	if got.Code != http.StatusConflict || replyOf(t, got).Error.Key != "statements.already_imported" {
		t.Errorf("second upload = %d %s, want 409 statements.already_imported", got.Code, got.Body.String())
	}
}

func TestUpload_NotMultipartIsAnInvalidRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/statements", strings.NewReader(`{"file":"x"}`))
	request.Header.Set("Content-Type", "application/json")

	got := newStatementsApp(t, testLimits, textExtractor{}).send(request)

	if got.Code != http.StatusBadRequest || replyOf(t, got).Error.Key != "common.invalid_request" {
		t.Errorf("upload = %d %s, want 400 common.invalid_request", got.Code, got.Body.String())
	}
}

func TestUpload_ReadingPastTheTimeoutIsCanceled(t *testing.T) {
	app := newStatementsApp(t, UploadLimits{MaxBytes: 1 << 20, Timeout: 50 * time.Millisecond}, textExtractor{delay: 200 * time.Millisecond})

	got := app.send(upload{files: [][]byte{[]byte(goodStatement)}, fields: withDate(nil)}.request(t))

	if got.Code == http.StatusCreated {
		t.Errorf("upload = %d, want an error once the timeout passed", got.Code)
	}
}

func TestUpload_ASlowClientIsCut(t *testing.T) {
	app := newStatementsApp(t, UploadLimits{MaxBytes: 1 << 20, Timeout: 100 * time.Millisecond}, textExtractor{})
	server := httptest.NewServer(app.router)
	t.Cleanup(server.Close)
	reader, writer := io.Pipe()
	t.Cleanup(func() { writer.Close() })
	go writer.Write([]byte("--x\r\nContent-Disposition: form-data; name=\"file\"; filename=\"c.pdf\"\r\n\r\nBANCO"))
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, server.URL+"/api/statements", reader)
	request.Header.Set("Content-Type", "multipart/form-data; boundary=x")
	started := time.Now()

	response, err := http.DefaultClient.Do(request)

	if err == nil {
		response.Body.Close()
	}
	if time.Since(started) > 3*time.Second || (err == nil && response.StatusCode == http.StatusCreated) {
		t.Errorf("slow upload took %v (err %v), want it cut soon after the timeout", time.Since(started), err)
	}
}

func TestBanks(t *testing.T) {
	got := newStatementsApp(t, testLimits, textExtractor{}).send(httptest.NewRequest(http.MethodGet, "/api/banks", nil))

	if got.Code != http.StatusOK || strings.TrimSpace(got.Body.String()) != `[{"id":"bci"}]` {
		t.Errorf("banks = %d %s, want 200 [{\"id\":\"bci\"}]", got.Code, got.Body.String())
	}
}
