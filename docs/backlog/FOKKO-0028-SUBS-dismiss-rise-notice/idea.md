# FOKKO-0028-SUBS · Descartar el aviso de alza

## Problema
El usuario revisa un alza y decide "mantener". En el MVP solo existe esa respuesta (`later` pasa al post-MVP) y descarta el aviso de ese monto.

## Depende de
- FOKKO-0026-SUBS y FOKKO-0021-SUBS.

## Alcance
- Puerto `PriceAlertDecisions` y caso de uso para descartar el aviso de un alza.
- Un alza descartada no vuelve a aparecer; una nueva alza por un monto mayor sí.

## Fuera
- `later` → post-MVP.
- Endpoint → FOKKO-0076-SUBS.

## Features afectadas
- `internal/subscriptions` (`ports`, `usecases`, `domain`)

## Dudas abiertas
- La tabla `price_alert_decisions` guarda una columna `decision` (`keep`/`later`). Con solo `keep` se propone no guardarla; se confirma al escribir FOKKO-0059-SUBS.

## Tareas
- [ ] [TASK-001 Puerto de decisiones de alza](TASK-001-PRICE-ALERT-PORT.md)
- [ ] [TASK-002 Descartar el aviso de un alza](TASK-002-KEEP-PRICE-ALERT.md)
- [ ] [TASK-003 El resumen oculta el alza descartada](TASK-003-SUMMARY-HIDES-KEPT-RISE.md)
