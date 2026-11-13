# FOKKO-0058A-SUBS/TASK-001 · Tabla `hidden_lines`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que una línea oculta siga oculta, para no verla de nuevo.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `UNIQUE(user_id, line_key)` y borrado en cascada con el usuario.

## Fuera
- Nada.

## Notas
- Ninguna.
