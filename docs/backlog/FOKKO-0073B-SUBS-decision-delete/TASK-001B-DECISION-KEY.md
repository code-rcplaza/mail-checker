# FOKKO-0073B-SUBS/TASK-001B · Validar la clave de línea
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dado `groupKey` ausente o vacío, cuando se envía, entonces responde 400 con `fields`.
- Dada una clave con `|`, `:`, `#` y espacios codificada en la URL, cuando se envía, entonces se decodifica idéntica.

## Fuera
- Nada.

## Notas
- Ninguna.
