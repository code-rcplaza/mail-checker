# FOKKO-0056-SUBS/TASK-002 · Adaptador de overrides
Feature: `internal/subscriptions/adapters/sqlite`

Como usuario quiero que mi override se use en el resumen, para ver mis números.

## Criterios
- Dado un override con solo monto, cuando se guarda y se lee, entonces el día queda vacío.
- Dados overrides de dos usuarios con la misma clave, cuando uno lista, entonces ve solo el suyo.
- Dado un override, cuando se borra, entonces ya no está; borrar uno inexistente no falla.

## Fuera
- Nada.

## Notas
- Ninguna.
