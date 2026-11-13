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
- La tabla se comparte entre dos features sin que una importe a la otra: cada una tiene su adaptador sobre la misma tabla. ¿Se acepta o va a `shared`?

## Tareas
- [ ] [TASK-001 Tabla `exchange_rates`](TASK-001-EXCHANGE-RATES-MIGRATION.md)
- [ ] [TASK-002 Guardar el tipo de cambio](TASK-002-EXCHANGE-RATE-WRITER.md)
- [ ] [TASK-003 Leer el más reciente](TASK-003-EXCHANGE-RATE-READER.md)
