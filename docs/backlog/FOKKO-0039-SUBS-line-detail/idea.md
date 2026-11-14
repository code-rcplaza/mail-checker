# FOKKO-0039-SUBS · `Summary`: detalle de una línea

## Problema
La pantalla Detalle muestra el ciclo, los últimos cobros y el alza de una línea.

## Depende de
- FOKKO-0023-SUBS y FOKKO-0026-SUBS.

## Alcance
- Datos de detalle de una línea: ciclo, últimos cobros y alza.

## Fuera
- Pantalla → FOKKO-0106A-SUBS y FOKKO-0106B-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Cuántos cobros recientes se listan: se propone los últimos 6.

## Tareas
- [x] [TASK-001 Ciclo de la línea](TASK-001-LINE-DETAIL-CYCLE.md)
- [x] [TASK-002 Últimos cobros](TASK-002-LINE-DETAIL-LAST-CHARGES.md)
- [x] [TASK-003 Detalle sin alza](TASK-003-LINE-DETAIL-WITHOUT-RISE.md)
