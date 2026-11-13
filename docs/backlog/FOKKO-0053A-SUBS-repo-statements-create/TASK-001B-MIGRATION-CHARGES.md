# FOKKO-0053A-SUBS/TASK-001B · Tabla `charges`
Feature: `internal/shared/sqlite/migrations`



## Criterios
- Dada `charges`, entonces `UNIQUE(user_id, source, external_id)`, `amount >= 0`, `currency IN ('CLP','USD')` y `card_last4` vacío o de 4 dígitos.
- Dado un usuario o una cartola borrados, cuando se consulta `charges`, entonces sus cargos desaparecieron (cascada).
- Dada la migración, cuando se aplica y se revierte, entonces funciona.

## Fuera
- Nada.

## Notas
- Ninguna.
