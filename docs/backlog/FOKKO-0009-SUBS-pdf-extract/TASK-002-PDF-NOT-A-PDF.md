# FOKKO-0009-SUBS/TASK-002 · Rechazar lo que no es un PDF
Feature: `internal/statements/adapters/extract`

Como usuario quiero un error claro si subo un archivo que no es PDF, para saber qué corregir.

## Criterios
- Dado un archivo que no empieza como PDF, cuando se extrae, entonces devuelve `ErrInvalid` con `statements.file_kind_unknown`.
- Dado un archivo vacío, cuando se extrae, entonces devuelve `ErrInvalid` con `statements.file_empty`.

## Fuera
- Nada.

## Notas
- Ninguna.
