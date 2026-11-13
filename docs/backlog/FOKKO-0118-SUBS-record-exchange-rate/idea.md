# FOKKO-0118-SUBS · Guardar el tipo de cambio derivado al importar una cartola

## Problema
La cartola trae el tipo de cambio, pero nadie lo guarda. Sin guardarlo, el total en CLP no puede convertir los dólares.

## Depende de
- FOKKO-0014-SUBS, FOKKO-0015-SUBS y FOKKO-0060-SUBS.

## Alcance
- `ImportStatement` deriva y guarda el tipo de cambio de la cartola, con su fecha.

## Fuera
- Leerlo para el total → FOKKO-0029-SUBS.

## Features afectadas
- `internal/statements` (`ports`, `usecases`)

## Dudas abiertas
- Dueño de la tabla y del puerto: `statements` lo escribe y `subscriptions` lo lee con su propio puerto `ExchangeRates`, sin importarse entre features. ¿Se acepta, o la tabla va a `shared`?

## Tareas
- [x] [TASK-001 Puerto para guardar el tipo de cambio](TASK-001-EXCHANGE-RATE-RECORDER-PORT.md)
- [x] [TASK-002 Importar guarda el tipo de cambio](TASK-002-IMPORT-RECORDS-RATE.md)
