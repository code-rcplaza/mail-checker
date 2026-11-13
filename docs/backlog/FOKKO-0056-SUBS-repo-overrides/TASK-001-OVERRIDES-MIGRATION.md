# FOKKO-0056-SUBS/TASK-001 · Tabla `line_overrides`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mi corrección de monto o día se conserve, para no repetirla.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `UNIQUE(user_id, line_key)`, `amount` y `billing_day` opcionales y `billing_day BETWEEN 1 AND 31`.
- Dado un día 0 o 32, cuando se inserta, entonces la base lo rechaza.

## Fuera
- Nada.

## Notas
- Ninguna.
