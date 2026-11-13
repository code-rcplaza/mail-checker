# FOKKO-0072A-SUBS/TASK-001 · Parámetro `months`
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero que se rechacen períodos absurdos, para no pedir 500 meses por error.

## Criterios
- Dado `months=12`, cuando se pide, entonces se usa 12.
- Dado `months` ausente, cuando se pide, entonces se usa el valor por defecto del caso de uso.
- Dado `months=0`, `-1`, `25` o `abc`, cuando se pide, entonces responde 400 con `fields`.

## Fuera
- Nada.

## Notas
- Ninguna.
