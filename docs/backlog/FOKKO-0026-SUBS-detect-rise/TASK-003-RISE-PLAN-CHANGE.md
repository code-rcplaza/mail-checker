# FOKKO-0026-SUBS/TASK-003 · Salto de más de 25 %
Feature: `internal/subscriptions/domain`

Como usuario quiero distinguir un cambio de plan de un reajuste, para entender el aumento.

## Criterios
- Dado un cobro de 20 a 30 (+50 %), cuando se detecta, entonces el alza viene marcada como cambio de plan.
- Dado un cobro de 20 a 25 (+25 %, justo el límite), cuando se detecta, entonces se documenta y prueba de qué lado cae.
- Dado un cobro de 100 a 100,5, cuando se detecta, entonces es alza sin marca de cambio de plan.

## Fuera
- Nada.

## Notas
- Ninguna.
