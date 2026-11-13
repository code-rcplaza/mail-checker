# FOKKO-0062A-SUBS/TASK-001 · Panic convertido en 500
Feature: `internal/shared/httpserver`

Como usuario quiero que un error inesperado no tire abajo la app, para seguir usándola.

## Criterios
- Dado un handler que hace panic, cuando llega una petición, entonces responde 500 con `common.internal` en JSON.
- Dado ese panic, cuando ocurre, entonces el detalle va al log y no a la respuesta.
- Dado el panic, cuando se atiende la petición siguiente, entonces el servidor sigue respondiendo.

## Fuera
- Nada.

## Notas
- Ninguna.
