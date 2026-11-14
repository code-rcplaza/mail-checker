package domain

import (
	"errors"
	"maps"
)

// Sentinels: the kind of failure. HTTP maps them to 400, 404, 409, 429 and 401.
var (
	ErrInvalid      = errors.New("invalid")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrTooMany      = errors.New("too many")
	ErrUnauthorized = errors.New("unauthorized")
)

// KeyedError is a domain failure that carries an i18n key for the front
// instead of a human message.
type KeyedError struct {
	Kind   error
	Key    MessageKey
	Params map[string]any
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
