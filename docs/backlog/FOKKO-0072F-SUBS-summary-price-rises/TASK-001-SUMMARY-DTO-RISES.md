# FOKKO-0072F-SUBS/TASK-001 · Alzas pendientes
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver qué suscripciones subieron y cuánto me cuesta, para decidir.

## Criterios
- Dada una alza, cuando se serializa, entonces trae `lineKey`, monto de antes y después, desde cuándo, impacto mensual y anual y si es cambio de plan.
- Dada una alza ya descartada, cuando se serializa, entonces no aparece.
- Dado un resumen sin alzas, cuando se serializa, entonces la lista viene vacía y no `null`.

## Fuera
- Nada.

## Notas
- Ninguna.
