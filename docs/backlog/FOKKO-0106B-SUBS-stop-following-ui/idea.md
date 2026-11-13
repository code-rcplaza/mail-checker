# FOKKO-0106B-SUBS · Dejar de seguir

## Problema
Detalle muestra una suscripción: ciclo, últimos cobros y el bloque de alza; desde ahí se edita, se cancela o se deja de seguir. El handoff solo diseñó el caso con alza.

Esta idea cubre solo: Dejar de seguir. El resto del trabajo original está en las ideas hermanas `FOKKO-0106*-SUBS`.

## Depende de
- FOKKO-0039-SUBS, FOKKO-0072A-SUBS a FOKKO-0072H-SUBS y FOKKO-0078-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Dejar de seguir

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
- [ ] [TASK-001A Pedir confirmación](TASK-001A-STOP-CONFIRM.md)
- [ ] [TASK-001B Accesibilidad del diálogo](TASK-001B-STOP-DIALOG-A11Y.md)
- [ ] [TASK-002 Dejar de seguir por HTTP](TASK-002-GATEWAY-STOP-FOLLOWING.md)
- [ ] [TASK-003 Conectar el hook con los componentes](TASK-003-CONTAINER-STOP.md)
