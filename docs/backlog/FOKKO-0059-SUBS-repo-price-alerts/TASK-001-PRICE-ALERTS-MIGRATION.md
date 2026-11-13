# FOKKO-0059-SUBS/TASK-001 · Tabla `price_alert_decisions`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que un alza descartada no vuelva, para no verla repetida.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `UNIQUE(user_id, line_key, amount)` y `amount > 0`.

## Fuera
- Nada.

## Notas
- Ninguna.
