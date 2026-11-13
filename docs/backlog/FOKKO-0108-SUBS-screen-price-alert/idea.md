# FOKKO-0108-SUBS · Alerta de alza

## Problema
Alerta avisa que una suscripción subió, con el monto de antes, el de después y el impacto. En el MVP solo se avisa cuando el cobro ya ocurrió; la respuesta es "mantener" o ir a cancelar.

## Depende de
- FOKKO-0076-SUBS y FOKKO-0027-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `PriceAlertCard` y `AlertDecisionActions`.

## Fuera
- "Recordármelo después": post-MVP.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 `PriceAlertCard`](TASK-001-PRICE-ALERT-CARD.md)
- [ ] [TASK-002 Mantener](TASK-002-ALERT-KEEP.md)
- [ ] [TASK-003 Ir a cancelar](TASK-003-ALERT-GO-CANCEL.md)
- [ ] [TASK-004 Enviar la decisión de alza por HTTP](TASK-004-GATEWAY-PRICE-ALERT.md)
- [ ] [TASK-005 Conectar el hook con los componentes](TASK-005-CONTAINER-PRICE-ALERT.md)
