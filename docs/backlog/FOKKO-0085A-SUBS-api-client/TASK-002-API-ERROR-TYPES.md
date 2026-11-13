# FOKKO-0085A-SUBS/TASK-002 · Errores de la API tipados
Feature: `frontend`

Como usuario quiero que cada pantalla sepa qué falló, para mostrar el mensaje correcto.

## Criterios
- Dado un error de la API, cuando se inspecciona, entonces expone `key`, `params`, `fields` y el status.
- Dado un error 401, cuando se inspecciona, entonces `isUnauthorized` es verdadero.
- Dado un error de red, cuando se inspecciona, entonces es un error distinto al de la API y se puede reintentar.

## Fuera
- Nada.

## Notas
- Ninguna.
