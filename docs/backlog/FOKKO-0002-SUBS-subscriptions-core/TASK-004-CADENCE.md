# FOKKO-0002-SUBS/TASK-004 · Cadencia
Feature: `internal/subscriptions/domain`

Como usuario quiero que se detecte si un cobro es mensual o anual aunque falten
algunos avisos.

## Criterios
- Dadas 3 fechas a 30 días, cuando se calcula la cadencia, entonces es mensual.
- Dadas fechas a 30, 61 y 30 días (un mes sin aviso), cuando se calcula, entonces es
  mensual.
- Dadas fechas con un intervalo de 120 días (más de 3 meses), cuando se calcula,
  entonces no hay cadencia.
- Dadas 2 fechas a 365 días, cuando se calcula, entonces es anual.
- Dada 1 sola fecha, o 2 fechas a 10 días, cuando se calcula, entonces no hay cadencia.
- Dadas fechas desordenadas, cuando se calcula, entonces el resultado es el mismo que
  ordenadas y la entrada no se modifica.

## Notas
- Mensual: cada intervalo es k × 24–38 días con k ≤ 3. Anual: 330–400 días.
  Constantes con nombre en `domain`.
