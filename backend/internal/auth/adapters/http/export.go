package http

import (
	"net/http"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/adapters/export"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/httpserver"
)

// copyFile is the downloadable copy in one format.
type copyFile struct {
	contentType string
	filename    string
	encode      func(export.UserCopy) ([]byte, error)
}

// copyFormats: CSV is the default, the one a spreadsheet opens.
var copyFormats = map[string]copyFile{
	"":     csvCopy,
	"csv":  csvCopy,
	"json": {contentType: "application/json", filename: "fokko-copia.json", encode: export.EncodeJSON},
}

var csvCopy = copyFile{contentType: "text/csv; charset=utf-8", filename: "fokko-copia.csv", encode: encodeCSVCopy}

func encodeCSVCopy(userCopy export.UserCopy) ([]byte, error) {
	return export.EncodeCSV(userCopy.Subscriptions)
}

func parseCopyFormat(raw string) (copyFile, error) {
	file, isKnown := copyFormats[raw]
	if !isKnown {
		return copyFile{}, shared.InvalidFields(map[string]shared.MessageKey{"format": domain.KeyExportFormatUnknown})
	}
	return file, nil
}

func (h Handlers) exportCopy(w http.ResponseWriter, r *http.Request) {
	exchange := httpserver.Exchange{Writer: w, Request: r}
	body, file, err := h.encodeCopy(r)
	if err != nil {
		h.Responder.WriteError(exchange, err)
		return
	}
	w.Header().Set("Content-Type", file.contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+file.filename+`"`)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		h.Responder.Logger.Warn("write copy failed", "err", err)
	}
}

// encodeCopy gathers the user's copy and encodes it in the asked format.
func (h Handlers) encodeCopy(r *http.Request) ([]byte, copyFile, error) {
	file, err := parseCopyFormat(r.URL.Query().Get("format"))
	if err != nil {
		return nil, copyFile{}, err
	}
	userID, err := httpserver.UserIDFrom(r.Context())
	if err != nil {
		return nil, copyFile{}, err
	}
	userCopy, err := h.Export.Gather(r.Context(), userID)
	if err != nil {
		return nil, copyFile{}, err
	}
	body, err := file.encode(userCopy)
	return body, file, err
}
