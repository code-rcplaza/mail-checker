# FOKKO-0059-SUBS · Repo de decisiones de alza

## Problema
Las alzas descartadas deben persistir con el monto que se descartó.

## Depende de
- FOKKO-0028-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `price_alert_decisions` y adaptador.

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Con `later` fuera del MVP no hace falta la columna `decision`.
- Decidido: se guarda el historial de montos mantenidos por línea (un monto que ya se mantuvo no vuelve a avisar); el resumen muestra el último monto cobrado. El puerto pasó a `map[GroupKey][]Money`.

## Tareas
- [x] [TASK-001 Tabla `price_alert_decisions`](TASK-001-PRICE-ALERTS-MIGRATION.md)
- [x] [TASK-002 Adaptador de alzas descartadas](TASK-002-PRICE-ALERTS-REPO.md)
