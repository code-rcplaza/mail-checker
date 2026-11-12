# 004-04 · `ImportStatement`: validaciones
Feature: `internal/statements/usecases`

Como usuario quiero mensajes claros cuando una cartola no se puede importar, sin que se
guarde nada a medias.

## Criterios
- Dado un archivo vacío, cuando se importa, entonces `ErrInvalid` con
  `statements.file_empty` y no se extrae.
- Dado un texto extraído vacío o solo espacios, cuando se importa, entonces `ErrInvalid`
  con `statements.file_empty`.
- Dada una fecha de facturación cero, cuando se importa, entonces `ErrInvalid` con
  `statements.billing_date_missing`.
- Dada una cartola sin ningún cobro válido, cuando se importa, entonces `ErrInvalid` con
  `statements.no_charges` y `badLines` = cantidad de líneas malas, y no se llama
  `Create`.
- Dado un `Create` que devuelve `ErrConflict` con `statements.already_imported`, cuando
  se importa, entonces el error se devuelve conservando su tipo y llave.
- Dado un error del extractor o del repositorio, cuando se importa, entonces se
  devuelve envuelto con contexto.

## Fuera
- `UNIQUE(bank, billing_date)` real en SQLite → fase 4.

## Notas
- El conflicto lo detecta el repositorio (atómico, sin carrera), no una consulta previa.
