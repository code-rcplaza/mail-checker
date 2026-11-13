# FOKKO-0007-SUBS/TASK-004 · Clave estable de la línea
Feature: `internal/subscriptions/domain`

Como usuario quiero que mis decisiones y alias sobre una línea se mantengan entre un
mes y otro.

## Criterios
- Dada una línea de Netflix con tarjeta `4821` cobrada los días 1, 2 y 1, cuando se
  calcula su clave, entonces `service:Netflix|4821|d01` (mediana del día, 2 dígitos).
- Dada una línea sin últimos 4 dígitos, cuando se calcula su clave, entonces la parte
  de tarjeta va vacía: `service:Netflix||d01`.
- Dadas dos líneas que calculan la misma clave, cuando se calculan, entonces la que
  empezó antes queda sin sufijo y la otra lleva `#2`.
- Dada la misma línea leída con 3 y con 6 meses de período (mismo día típico), cuando
  se calcula su clave, entonces es la misma.
- Dado un padre de comercio desconocido, cuando se calcula la clave de su línea,
  entonces conserva el padre: `merchant:SEGUROS CLC S A:CLP|4821|d04`.
- Ninguna clave contiene ids internos ni correos.

## Fuera
- Uso de la clave en `GroupCharges` → FOKKO-0007-SUBS/TASK-005.

## Notas
- Con un número par de cobros, la mediana es el menor de los dos días centrales.
