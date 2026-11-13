# FOKKO-0055-SUBS/TASK-001 · Tabla `line_aliases`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mis nombres de línea queden guardados, para no repetirlos.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `UNIQUE(user_id, line_key)` y el alias no pasa de 40 caracteres.
- Dado un alias de 41 caracteres, cuando se inserta, entonces la base lo rechaza.

## Fuera
- Nada.

## Notas
- Ninguna.
