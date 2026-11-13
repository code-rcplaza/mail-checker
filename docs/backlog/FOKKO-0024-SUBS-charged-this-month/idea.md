# FOKKO-0024-SUBS · Cobros ya cobrados del mes

## Problema
"Llevas 3 cobros este mes" necesita contar las suscripciones ya cobradas en el mes calendario actual.

## Depende de
- FOKKO-0023-SUBS (usa el reloj y el mes en `America/Santiago`).

## Alcance
- Contar las suscripciones con un cobro en el mes actual.

## Fuera
- Los restantes → FOKKO-0025-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Contar los cobrados del mes](TASK-001-COUNT-CHARGED-THIS-MONTH.md)
- [ ] [TASK-002 El resumen trae los cobrados](TASK-002-SUMMARY-CHARGED-COUNT.md)
