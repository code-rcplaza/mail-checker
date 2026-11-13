# FOKKO-0004-SUBS/TASK-005 · `ListStatements`, `RemoveStatement`, `ListBanks`
Feature: `internal/statements/usecases`

Como usuario quiero ver mis cartolas, borrar una y saber qué bancos están soportados.

## Criterios
- Dadas cartolas guardadas, cuando se llama `ListStatements`, entonces devuelve las del
  repositorio.
- Dado un `publicID` existente, cuando se llama `RemoveStatement`, entonces se llama
  `Delete` con ese `publicID`.
- Dado un `publicID` que no existe, cuando se llama `RemoveStatement`, entonces
  `ErrNotFound` con `statements.not_found`.
- Dados parsers registrados de varios bancos, cuando se llama `ListBanks`, entonces
  devuelve sus bancos ordenados alfabéticamente.
- Dado un error del repositorio, cuando se lista o se borra, entonces se devuelve
  envuelto con contexto.

## Fuera
- Endpoints `/api/statements` y `/api/banks` → fase 5.
