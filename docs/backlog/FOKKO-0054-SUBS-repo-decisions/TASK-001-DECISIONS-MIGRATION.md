# FOKKO-0054-SUBS/TASK-001 · Tabla `subscription_decisions`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mis decisiones sobrevivan, para no responder lo mismo cada vez.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `UNIQUE(user_id, group_key)` y `decision IN ('confirmed','dismissed')`.
- Dada una decisión confirmada sin cadencia, cuando se inserta, entonces la base la rechaza; una descartada con cadencia también.

## Fuera
- Nada.

## Notas
- Ninguna.
