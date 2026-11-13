# FOKKO-0029-SUBS · Total del mes en CLP

## Problema
El total del mes va en CLP: los cobros en USD se convierten con el último tipo de cambio de la cartola (decisión 5 y 6). Si no hay tipo de cambio, el total avisa que está incompleto.

## Depende de
- FOKKO-0016-SUBS y FOKKO-0017-SUBS.

## Alcance
- `Summary` trae el total del mes en CLP y si está incompleto.
- Puerto de solo lectura `ExchangeRates` (último tipo de cambio por usuario).

## Fuera
- Guardar el tipo de cambio → FOKKO-0118-SUBS.
- Mostrar el monto original → FOKKO-0030-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`, `ports`, `usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Sumar el gasto convertido a pesos](TASK-001A-TOTAL-SUM.md)
- [x] [TASK-001B Marcar el total aproximado](TASK-001B-TOTAL-APPROX-FLAG.md)
- [x] [TASK-002 Total sin tipo de cambio](TASK-002-TOTAL-INCOMPLETE.md)
- [x] [TASK-003 `GetSummary` lee el tipo de cambio](TASK-003-SUMMARY-READS-RATE.md)
