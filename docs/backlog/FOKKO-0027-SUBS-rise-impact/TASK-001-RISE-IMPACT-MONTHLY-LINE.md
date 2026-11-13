# FOKKO-0027-SUBS/TASK-001 · Impacto en una línea mensual
Feature: `internal/subscriptions/domain`

Como usuario quiero ver cuánto más pago al año por esta alza, para decidir si sigo.

## Criterios
- Dada una alza de 10.000 a 10.500 en línea mensual, cuando se calcula, entonces el impacto anual es 6.000 y el mensual 500.
- Dado un impacto con decimales, cuando se calcula, entonces redondea half-up a centésimas.

## Fuera
- Nada.

## Notas
- Ninguna.
