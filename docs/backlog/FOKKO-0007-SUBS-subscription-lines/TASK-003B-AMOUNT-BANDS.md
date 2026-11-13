# FOKKO-0007-SUBS/TASK-003B · Líneas por banda de monto
Feature: `internal/subscriptions/domain`

Como usuario quiero que dos suscripciones de montos distintos que llegan al mismo
correo nunca se mezclen, aunque se cobren en días cercanos al cambio de mes.

## Criterios
- Dados $12.990 los días 1 y $7.990 los días 28 de agosto, septiembre y octubre (mismo
  origen), cuando se encadenan, entonces quedan dos líneas de tres cobros, cada una
  con un solo monto (sin cambio de plan falso).
- Dados $12.990, $12.990 y luego $17.990 a un mes (mismo origen), cuando se encadenan,
  entonces sigue siendo una sola línea (las bandas que se suceden se unen).
- Dados CLP 5.990, CLP 5.990 y luego USD 7,99, USD 7,99 a un mes cada uno, cuando se
  encadenan, entonces es una sola línea (cambio de plan con cambio de moneda).
- Dados $12.990 mensuales de agosto a octubre y $7.990 mensuales de septiembre a
  octubre (conviven en el tiempo), cuando se encadenan, entonces son dos líneas.
- Los criterios de FOKKO-0007-SUBS/TASK-003 siguen pasando sin cambios.

## Fuera
- Clave de la línea → FOKKO-0007-SUBS/TASK-004.

## Notas
- Pasos: agrupar los cobros por banda de monto (misma moneda, ±25 % del último cobro
  de la banda); encadenar dentro de cada banda (regla de FOKKO-0007-SUBS/TASK-003); unir una línea con
  otra que empieza a un intervalo válido después de que la primera terminó (cambio de
  plan); separar los cobros sueltos.
- Límite aceptado: dos suscripciones del **mismo** monto cobradas en días a menos de
  ~4 días del cambio de mes pueden intercambiar su primer cobro; no hay cambio de plan
  falso ni cambia el gasto.
