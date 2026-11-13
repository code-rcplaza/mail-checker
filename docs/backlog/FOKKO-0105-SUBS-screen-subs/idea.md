# FOKKO-0105-SUBS · Subs

## Problema
Subs lista todas las suscripciones del usuario, con su monto, próximo cobro y marcas.

## Depende de
- FOKKO-0072A-SUBS a FOKKO-0072H-SUBS y FOKKO-0088-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `SubscriptionList` y `SubscriptionRow`.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 `SubscriptionList`](TASK-001-SUBSCRIPTION-LIST.md)
- [ ] [TASK-002A Datos básicos de la fila](TASK-002A-ROW-BASIC.md)
- [ ] [TASK-002B Monto en dólares](TASK-002B-ROW-USD.md)
- [ ] [TASK-002C Marcas de la fila](TASK-002C-ROW-BADGES.md)
- [ ] [TASK-003 Subs vacío](TASK-003-SUBS-EMPTY.md)
- [ ] [TASK-004 Leer las suscripciones](TASK-004-HOOK-SUBSCRIPTIONS.md)
- [ ] [TASK-005 Leer las suscripciones por HTTP](TASK-005-GATEWAY-SUBSCRIPTIONS.md)
- [ ] [TASK-006 Conectar el hook con los componentes](TASK-006-CONTAINER-SUBS.md)
