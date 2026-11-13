# FOKKO-0037-SUBS/TASK-001 · Acciones pendientes
Feature: `internal/subscriptions/domain`

Como usuario quiero ver qué me falta resolver, para dejar mi lista al día.

## Criterios
- Dadas 2 posibles vigentes y 1 alza, cuando se arma, entonces hay 3 acciones pendientes.
- Dada una posible fuera de su ventana (inactiva), cuando se arma, entonces no pide acción.
- Dada un alza ya descartada, cuando se arma, entonces no pide acción.

## Fuera
- Nada.

## Notas
- Ninguna.
