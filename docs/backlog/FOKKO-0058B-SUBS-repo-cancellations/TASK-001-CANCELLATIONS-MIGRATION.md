# FOKKO-0058B-SUBS/TASK-001 · Tabla `cancellations`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mi marca "por cancelar" se conserve, para recordarla.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `UNIQUE(user_id, line_key)`, guarda la fecha de la marca y se borra en cascada con el usuario.

## Fuera
- Nada.

## Notas
- Ninguna.
