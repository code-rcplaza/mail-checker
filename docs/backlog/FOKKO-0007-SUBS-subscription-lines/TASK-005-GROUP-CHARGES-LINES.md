# FOKKO-0007-SUBS/TASK-005 · `GroupCharges` produce líneas
Feature: `internal/subscriptions/domain`

Como usuario quiero que la agrupación entregue una línea por suscripción, con su padre.

## Criterios
- Dados cobros `APPLE.COM/BILL` de $4.690 y de $5.990, cada uno mensual, cuando se
  agrupan, entonces hay un padre `umbrella:App Store` con dos líneas (ya no un grupo
  por monto exacto).
- Dados los cobros de dos Netflix de la misma cuenta, cuando se agrupan, entonces hay
  dos `ChargeGroup` con el mismo padre `service:Netflix` y claves de línea distintas.
- Dado un cobro suelto que FOKKO-0007-SUBS/TASK-003 deja fuera de la línea, cuando se agrupa, entonces
  queda en los otros cobros de su padre, no en una línea.
- Los grupos salen ordenados por clave y `GroupCharges` no modifica el slice recibido.

## Fuera
- Resumen y gasto → FOKKO-0007-SUBS/TASK-006.

## Notas
- ⚠️ Cambia el diseño de FOKKO-0002-SUBS/TASK-002 y FOKKO-0002-SUBS/TASK-012: hay tests existentes que esperan claves
  `umbrella:<nombre>:<monto>:<moneda>` y "otros cobros del servicio" dentro del grupo.
  Antes de tocarlos se pide el ok del usuario (regla de tests intocables), mostrando
  cada test y su nuevo `want`.
