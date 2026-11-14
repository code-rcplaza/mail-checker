# FOKKO-0060-SUBS · Repo del tipo de cambio

## Problema
El tipo de cambio derivado de cada cartola debe guardarse con su fecha, y `subscriptions` debe leer el más reciente.

## Depende de
- FOKKO-0016-SUBS, FOKKO-0049-SUBS y FOKKO-0118-SUBS.

## Alcance
- Migración de `exchange_rates`, escritor en `statements` y lector en `subscriptions`.

## Fuera
- Nada.

## Features afectadas
- `internal/statements` (`adapters/sqlite`)
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Resuelto como `charges` (FOKKO-0053C-SUBS): cada feature tiene su adaptador sobre la misma tabla. "Leer el último" lo hace `LatestExchangeRate` en el dominio; el lector entrega los del usuario ordenados por fecha.

## Tareas
- [x] [TASK-001 Tabla `exchange_rates`](TASK-001-EXCHANGE-RATES-MIGRATION.md)
- [x] [TASK-002 Guardar el tipo de cambio](TASK-002-EXCHANGE-RATE-WRITER.md)
- [x] [TASK-003 Leer el más reciente](TASK-003-EXCHANGE-RATE-READER.md)
