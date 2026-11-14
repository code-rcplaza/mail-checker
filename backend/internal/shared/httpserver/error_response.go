package httpserver

import "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"

// errorResponse is the only error shape the front receives: an i18n key, never text.
type errorResponse struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Key    domain.MessageKey            `json:"key"`
	Params map[string]any               `json:"params,omitempty"`
	Fields map[string]domain.MessageKey `json:"fields,omitempty"`
}
