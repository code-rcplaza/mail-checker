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
- Con `later` fuera del MVP no hace falta la columna `decision`; la tabla guarda solo el descarte. Se actualiza `ARCHITECTURE.md` cuando se escriba.

## Tareas
- [ ] [TASK-001 Tabla `price_alert_decisions`](TASK-001-PRICE-ALERTS-MIGRATION.md)
- [ ] [TASK-002 Adaptador de alzas descartadas](TASK-002-PRICE-ALERTS-REPO.md)
