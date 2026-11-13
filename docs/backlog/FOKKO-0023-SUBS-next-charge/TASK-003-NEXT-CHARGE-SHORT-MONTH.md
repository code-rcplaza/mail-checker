# FOKKO-0023-SUBS/TASK-003 · Día 31 en meses cortos
Feature: `internal/subscriptions/domain`

Como usuario quiero que un cobro del día 31 no desaparezca en abril, para ver una fecha real.

## Criterios
- Dado día típico 31 y el mes siguiente de 30 días, cuando se calcula, entonces el próximo cae en el día 30.
- Dado día típico 30 y febrero, cuando se calcula, entonces cae en el último día de febrero.

## Fuera
- Nada.

## Notas
- Ninguna.
