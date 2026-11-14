// Package http is the statements feature's HTTP edge: the upload and the bank list.
package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/statements/usecases"
)

// UploadLimits bound a statement upload: its size and how long reading and
// importing it may take.
type UploadLimits struct {
	MaxBytes int64
	Timeout  time.Duration
}

// DefaultUploadLimits: a statement PDF weighs a few MB.
var DefaultUploadLimits = UploadLimits{MaxBytes: 10 << 20, Timeout: 60 * time.Second}

type Handlers struct {
	Import    usecases.ImportStatement
	Banks     usecases.ListBanks
	Limits    UploadLimits
	Responder httpserver.Responder
}

// MountProtected adds the routes that go behind the session.
func (h Handlers) MountProtected(router chi.Router) {
	router.Post("/api/statements", h.upload)
	router.Get("/api/banks", h.banks)
}

type uploadResponse struct {
	ID       string   `json:"id"`
	Filename string   `json:"filename"`
	Months   []string `json:"months"`
	Charges  int      `json:"charges"`
	BadLines int      `json:"badLines"` // how many; their text never leaves the back
}

type bankResponse struct {
	ID shared.Bank `json:"id"`
}

func (h Handlers) upload(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	ctx, cancel := h.limitUpload(exchange)
	defer cancel()
	userID, err := httpserver.UserIDFrom(ctx)
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	input, err := h.readUpload(exchange)
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	result, err := h.Import.Execute(ctx, userID, input)
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	h.Responder.Respond(exchange, httpserver.Reply{Status: http.StatusCreated, Body: uploadResponseOf(result)})
}

// limitUpload cuts a client that sends too slowly and gives the import a
// deadline; the body is capped so an oversized file is never read whole.
func (h Handlers) limitUpload(exchange httpserver.Exchange) (context.Context, context.CancelFunc) {
	deadline := time.Now().Add(h.Limits.Timeout)
	if err := http.NewResponseController(exchange.Writer).SetReadDeadline(deadline); err != nil && !errors.Is(err, http.ErrNotSupported) {
		h.Responder.Logger.Warn("set upload read deadline failed", "err", err)
	}
	exchange.Request.Body = http.MaxBytesReader(exchange.Writer, exchange.Request.Body, h.Limits.MaxBytes)
	return context.WithDeadline(exchange.Request.Context(), deadline)
}

func (h Handlers) readUpload(exchange httpserver.Exchange) (usecases.ImportStatementInput, error) {
	request := exchange.Request
	if err := requireMultipart(request); err != nil {
		return usecases.ImportStatementInput{}, err
	}
	if err := parseForm(request, h.Limits.MaxBytes); err != nil {
		return usecases.ImportStatementInput{}, err
	}
	defer request.MultipartForm.RemoveAll()
	file, err := onlyFile(request.MultipartForm)
	if err != nil {
		return usecases.ImportStatementInput{}, err
	}
	billingDate, err := parseBillingDate(request.FormValue("billingDate"))
	if err != nil {
		return usecases.ImportStatementInput{}, err
	}
	data, err := readFile(file)
	if err != nil {
		return usecases.ImportStatementInput{}, err
	}
	return usecases.ImportStatementInput{Filename: file.Filename, Data: data, BillingDate: billingDate, Bank: shared.Bank(request.FormValue("bank"))}, nil
}

func requireMultipart(request *http.Request) error {
	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" {
		return shared.Invalid(shared.KeyInvalidRequest, nil)
	}
	return nil
}

func parseForm(request *http.Request, maxBytes int64) error {
	err := request.ParseMultipartForm(maxBytes)
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return shared.Invalid(domain.KeyFileTooLarge, map[string]any{"maxBytes": maxBytes})
	}
	if err != nil {
		return fmt.Errorf("read upload: %v: %w", err, shared.Invalid(shared.KeyInvalidRequest, nil))
	}
	return nil
}

func onlyFile(form *multipart.Form) (*multipart.FileHeader, error) {
	files := form.File["file"]
	switch len(files) {
	case 0:
		return nil, shared.InvalidFields(map[string]shared.MessageKey{"file": domain.KeyFileRequired})
	case 1:
		return files[0], nil
	default:
		return nil, shared.Invalid(domain.KeyFileTooMany, nil)
	}
}

// parseBillingDate reads unix seconds, like every date in the API; empty means not given.
func parseBillingDate(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}
	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, shared.InvalidFields(map[string]shared.MessageKey{"billingDate": domain.KeyBillingDateInvalid})
	}
	return time.Unix(seconds, 0).UTC(), nil
}

func readFile(header *multipart.FileHeader) ([]byte, error) {
	file, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read uploaded file: %w", err)
	}
	return data, nil
}

func uploadResponseOf(result usecases.ImportResult) uploadResponse {
	months := make([]string, 0, len(result.Months))
	for _, month := range result.Months {
		months = append(months, fmt.Sprintf("%04d-%02d", month.Year, month.Month))
	}
	statement := result.Statement
	return uploadResponse{ID: statement.PublicID, Filename: statement.Filename, Months: months, Charges: statement.Charges, BadLines: len(statement.BadLines)}
}

func (h Handlers) banks(w http.ResponseWriter, r *http.Request) {
	banks := h.Banks.Execute()
	response := make([]bankResponse, 0, len(banks))
	for _, bank := range banks {
		response = append(response, bankResponse{ID: bank})
	}
	h.Responder.Respond(httpserver.Exchange{Writer: w, Request: r}, httpserver.Reply{Status: http.StatusOK, Body: response})
}
