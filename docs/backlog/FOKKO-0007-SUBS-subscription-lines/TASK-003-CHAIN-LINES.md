# FOKKO-0007-SUBS/TASK-003 · Encadenar líneas
Feature: `internal/subscriptions/domain`

Como usuario quiero que cada suscripción de un mismo servicio quede en su propia línea,
aunque compartan cuenta y tarjeta.

## Criterios
- Dados cobros de $12.990 los días 25 y $7.990 los días 9 de tres meses (mismo
  origen), cuando se encadenan, entonces quedan dos líneas de tres cobros cada una.
- Dados dos cobros de $5.990 por mes, los días 5 y 9, durante tres meses (mismo
  origen), cuando se encadenan, entonces quedan dos líneas de tres cobros.
- Dados dos cobros por mes el mismo día y del mismo monto, cuando se encadenan,
  entonces quedan dos líneas de tres cobros (una por cada cobro del mes).
- Dados cobros de $12.990, $12.990 y luego $17.990 en la fecha esperada (mismo
  origen), cuando se encadenan, entonces es una sola línea (cambio de plan).
- Dado un cobro suelto de USD 5 entre cobros mensuales de USD 20 del mismo origen,
  cuando se encadenan, entonces la línea tiene los de USD 20 y el de USD 5 queda como
  otro cobro del padre.
- Dado un único cobro en un origen sin otras líneas, cuando se encadena, entonces es
  una línea de un cobro.
- Dados cobros anuales (día 0 y día 365), cuando se encadenan, entonces es una línea.

## Fuera
- Clave de la línea → FOKKO-0007-SUBS/TASK-004.

## Notas
- Regla: un cobro continúa la línea si está a un intervalo válido de su último cobro
  (mensual k × 24–38 o anual 330–400); con varias candidatas, monto más parecido y
  luego fecha esperada más cercana; si ninguna, línea nueva. Un salto > 25 % en la
  fecha de la línea es cambio de plan. Un cobro que no encadena, en un origen con una
  línea activa, es otro cobro del padre.
