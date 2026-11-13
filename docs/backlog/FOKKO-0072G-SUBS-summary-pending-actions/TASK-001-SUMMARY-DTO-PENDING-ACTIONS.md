# FOKKO-0072G-SUBS/TASK-001 · Acciones pendientes
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver qué me falta resolver, para dejar mi lista al día.

## Criterios
- Dadas alzas y posibles pendientes, cuando se serializa, entonces vienen en el orden del dominio con su tipo y `lineKey`.
- Dado ningún pendiente, cuando se serializa, entonces la lista viene vacía y no `null`.

## Fuera
- Nada.

## Notas
- Ninguna.
