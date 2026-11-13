# FOKKO-0103-SUBS · Resumen

## Problema
El Resumen es la pantalla principal: total del mes, reparto por servicio, próximo cobro, acciones pendientes y próximos cobros.

## Depende de
- FOKKO-0072A-SUBS a FOKKO-0072H-SUBS, FOKKO-0089A-SUBS y FOKKO-0089B-SUBS y FOKKO-0092-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `SpendHero`, `ServiceLegend`, `NextChargeRow`, `PendingActionsList`, `UpcomingChargesList` y la línea de cobros del mes.

## Fuera
- `CycleTip` → FOKKO-0104-SUBS.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 `SpendHero`](TASK-001-SPEND-HERO.md)
- [ ] [TASK-002 `ServiceLegend`](TASK-002-SERVICE-LEGEND.md)
- [ ] [TASK-003 Cobros del mes](TASK-003-MONTH-COUNTS.md)
- [ ] [TASK-004 `NextChargeRow`](TASK-004-NEXT-CHARGE-ROW.md)
- [ ] [TASK-005 `PendingActionsList`](TASK-005-PENDING-ACTIONS-LIST.md)
- [ ] [TASK-006 `UpcomingChargesList`](TASK-006-UPCOMING-CHARGES-LIST.md)
- [ ] [TASK-007 Resumen vacío](TASK-007-SUMMARY-EMPTY.md)
- [ ] [TASK-008 Leer el resumen](TASK-008-HOOK-SUMMARY.md)
- [ ] [TASK-009 Leer el resumen por HTTP](TASK-009-GATEWAY-SUMMARY.md)
- [ ] [TASK-010 Conectar el hook con los componentes](TASK-010-CONTAINER-SUMMARY.md)
