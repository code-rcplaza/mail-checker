# 002-03 · Regla de coincidencia cartola ↔ correo
Feature: `internal/subscriptions/domain`

Como usuario quiero que un cobro que aparece en la cartola y también llegó por
correo cuente una sola vez.

## Criterios
- Dados un cobro de cartola y uno de correo con mismo banco, servicio y monto, a 20
  días, cuando se comparan, entonces coinciden.
- Dados los mismos cobros a 36 días, cuando se comparan, entonces no coinciden.
- Dados cobros con distinto banco, o distinto monto, o distinto servicio, cuando se
  comparan, entonces no coinciden.
- Dada una lista con un cobro de cartola que coincide con uno de correo y otro que no,
  cuando se quitan los duplicados, entonces queda el de correo y el de cartola sin
  pareja.
- Dados dos cobros de cartola iguales y un solo cobro de correo, cuando se quitan los
  duplicados, entonces solo uno de cartola se descarta (emparejamiento 1 a 1).

## Fuera
- La vista diferenciada → 002-10.

## Notas
- "Mismo servicio" usa la misma identidad que la agrupación (servicio conocido o
  nombre normalizado).
