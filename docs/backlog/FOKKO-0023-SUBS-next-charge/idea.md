# FOKKO-0023-SUBS · Próximo cobro por línea

## Problema
"Próximo cobro" es lo primero que quiere ver el usuario. Hoy el resumen no lo calcula.

## Alcance
- Función pura: el último cobro más un intervalo de la cadencia (mes o año), en el día típico de la línea.
- `Subscription` lo trae en el resumen.

## Fuera
- Próximo cobro de una manual → FOKKO-0031-SUBS.
- Próximo cobro con override → FOKKO-0033-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Si la fecha calculada ya pasó (cobro atrasado hasta 45 días), se muestra tal cual, sin saltar al siguiente. ¿Se acepta?

## Tareas
- [ ] [TASK-001 Próximo cobro mensual](TASK-001-NEXT-CHARGE-MONTHLY.md)
- [ ] [TASK-002 Próximo cobro anual](TASK-002-NEXT-CHARGE-YEARLY.md)
- [ ] [TASK-003 Día 31 en meses cortos](TASK-003-NEXT-CHARGE-SHORT-MONTH.md)
- [ ] [TASK-004 `Subscription` trae su próximo cobro](TASK-004-SUMMARY-NEXT-CHARGE.md)
