package domain

import (
	"errors"
	"maps"
)

// Sentinels: the kind of failure. HTTP maps them to 400, 404, 409, 429, 401, 403 and 503.
var (
	ErrInvalid      = errors.New("invalid")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrTooMany      = errors.New("too many")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrUnavailable  = errors.New("unavailable")
)

// KeyedError is a domain failure that carries an i18n key for the front
// instead of a human message.
type KeyedError struct {
	Kind   error
	Key    MessageKey
	Params map[string]any
	Fields map[string]MessageKey
}

// Error is for logs only; users see the translated Key.
func (e KeyedError) Error() string { return e.Kind.Error() + ": " + string(e.Key) }

func (e KeyedError) Unwrap() error { return e.Kind }

func Invalid(key MessageKey, params map[string]any) error {
	return KeyedError{Kind: ErrInvalid, Key: key, Params: maps.Clone(params)}
}

func NotFound(key MessageKey) error {
	return KeyedError{Kind: ErrNotFound, Key: key}
}

func Conflict(key MessageKey) error {
	return KeyedError{Kind: ErrConflict, Key: key}
}

func TooMany(key MessageKey, params map[string]any) error {
	return KeyedError{Kind: ErrTooMany, Key: key, Params: maps.Clone(params)}
}

func Unauthorized(key MessageKey) error {
	return KeyedError{Kind: ErrUnauthorized, Key: key}
}

// InvalidFields is a validation failure with the key of each wrong field, so
// the front shows every error next to its field.
func InvalidFields(fields map[string]MessageKey) error {
	return KeyedError{Kind: ErrInvalid, Key: KeyValidationFailed, Fields: maps.Clone(fields)}
}

func Forbidden(key MessageKey) error {
	return KeyedError{Kind: ErrForbidden, Key: key}
}

func Unavailable(key MessageKey) error {
	return KeyedError{Kind: ErrUnavailable, Key: key}
}

// FieldErrors joins the validation of each field into one InvalidFields, or
// nil when every field passed. An error without key counts as a plain failure.
func FieldErrors(fields map[string]error) error {
	failed := make(map[string]MessageKey)
	for field, err := range fields {
		if err != nil {
			failed[field] = keyOf(err)
		}
	}
	if len(failed) == 0 {
		return nil
	}
	return InvalidFields(failed)
}

func keyOf(err error) MessageKey {
	var keyed KeyedError
	if errors.As(err, &keyed) {
		return keyed.Key
	}
	return KeyValidationFailed
}
