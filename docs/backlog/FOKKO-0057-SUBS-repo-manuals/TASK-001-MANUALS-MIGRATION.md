# FOKKO-0057-SUBS/TASK-001 · Tabla `manual_subscriptions`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mis suscripciones manuales se guarden, para verlas siempre.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces `public_id` único, `billing_day BETWEEN 1 AND 31`, `amount > 0` y moneda en CLP o USD.
- Dado un usuario borrado, cuando se consulta, entonces sus manuales desaparecieron.

## Fuera
- Nada.

## Notas
- Ninguna.
