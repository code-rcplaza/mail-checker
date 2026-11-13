# FOKKO-0070-SUBS/TASK-002A · Errores 400 de la subida
Feature: `internal/statements/adapters/http`



## Criterios
- Dado un banco no detectado, cuando se sube, entonces responde 400 con `statements.bank_not_detected`.
- Dado un archivo sin cargos, cuando se sube, entonces responde 400 con `statements.no_charges` y `badLines`.
- Dado un archivo que no es PDF, cuando se sube, entonces responde 400 con `statements.file_kind_unknown`.

## Fuera
- Nada.

## Notas
- Ninguna.
