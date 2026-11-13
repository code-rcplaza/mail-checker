# FOKKO-0001-SUBS/TASK-007 · Mapeo de errores a HTTP
Feature: `internal/shared/httpserver`

Como usuario quiero que el front reciba siempre una llave i18n para mostrarme un
mensaje claro, sin detalles internos.

## Criterios
- Dado `Invalid("statements.bank_not_detected", {"count": 2})`, cuando se responde,
  entonces 400 y `{"error":{"key":"statements.bank_not_detected","params":{"count":2}}}`.
- Dado `NotFound(key)` o `Conflict(key)`, cuando se responde, entonces 404 o 409 con
  esa llave.
- Dado `ErrNotFound` o `ErrConflict` sin llave propia, cuando se responde, entonces
  404 con `common.not_found` o 409 con `common.conflict`.
- Dado un error de dominio envuelto (`fmt.Errorf("...: %w", err)`), cuando se
  responde, entonces el status y la llave son los mismos.
- Dado cualquier otro error, cuando se responde, entonces 500 con
  `{"error":{"key":"common.internal"}}`, sin el mensaje original, y el error real
  queda en el log.
- Dado un error sin `params`, cuando se responde, entonces el JSON no incluye
  `params`.

## Fuera
- Llaves específicas de cada feature (constantes en su `domain`, fases siguientes).
