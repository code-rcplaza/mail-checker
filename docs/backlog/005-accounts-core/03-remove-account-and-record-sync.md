# 005-03 · `RemoveAccount` y `RecordSync`
Feature: `internal/accounts/usecases`

Como usuario quiero borrar una cuenta y ver cuándo se sincronizó por última vez cada
una.

## Criterios
- Dado un `publicID` existente, cuando se llama `RemoveAccount`, entonces se llama
  `Delete` con ese `publicID`.
- Dado un `publicID` que no existe, cuando se llama `RemoveAccount`, entonces
  `ErrNotFound` con `accounts.not_found`.
- Dados un id, una hora y una llave de error, cuando se llama `RecordSync`, entonces se
  llama `MarkSynced` con esos valores (llave vacía = sin error).
- Dada una cuenta borrada durante la sincronización (`MarkSynced` devuelve
  `ErrNotFound`), cuando se llama `RecordSync`, entonces devuelve nil.
- Dado otro error del repositorio, cuando se borra o se registra, entonces se devuelve
  envuelto con contexto.

## Fuera
- Borrado en cascada de los cobros → repositorio SQLite (fase 4).

## Notas
- Ignorar la cuenta borrada evita que una cuenta corte la sincronización de las demás.
