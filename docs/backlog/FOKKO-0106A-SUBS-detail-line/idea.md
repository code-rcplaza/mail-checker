# FOKKO-0106A-SUBS · Detalle de una línea (con o sin alza)

## Problema
Detalle muestra una suscripción: ciclo, últimos cobros y el bloque de alza; desde ahí se edita, se cancela o se deja de seguir. El handoff solo diseñó el caso con alza.

Esta idea cubre solo: Detalle de una línea (con o sin alza). El resto del trabajo original está en las ideas hermanas `FOKKO-0106*-SUBS`.

## Depende de
- FOKKO-0039-SUBS, FOKKO-0072A-SUBS a FOKKO-0072H-SUBS y FOKKO-0078-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Encabezado y ciclo
- `PriceChangeBlock`
- `ChargeHistory`

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- Falta el diseño de Detalle sin alza; se propone el mismo bloque sin la tarjeta de alza.
- La confirmación de "dejar de seguir" es un diálogo: hay que cuidar foco, Esc y bloqueo de scroll a mano. Si resulta frágil se propone Radix Primitives como dependencia.

## Tareas
- [ ] [TASK-001A Ciclo mensual](TASK-001A-DETAIL-MONTHLY.md)
- [ ] [TASK-001B Ciclo anual](TASK-001B-DETAIL-YEARLY.md)
- [ ] [TASK-001C Alias en el encabezado](TASK-001C-DETAIL-ALIAS.md)
- [ ] [TASK-002 `PriceChangeBlock`](TASK-002-PRICE-CHANGE-BLOCK.md)
- [ ] [TASK-003 `ChargeHistory`](TASK-003-CHARGE-HISTORY.md)
- [ ] [TASK-004 Leer el detalle de una línea](TASK-004-HOOK-LINE-DETAIL.md)
- [ ] [TASK-005 Leer el detalle por HTTP](TASK-005-GATEWAY-LINE-DETAIL.md)
- [ ] [TASK-006 Conectar el hook con los componentes](TASK-006-CONTAINER-DETAIL.md)
