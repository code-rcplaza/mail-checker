# FOKKO-0070-SUBS/TASK-001 · Subir una cartola
Feature: `internal/statements/adapters/http`

Como usuario quiero subir mi cartola y ver qué se leyó, para seguir al análisis.

## Criterios
- Dado un PDF válido, cuando se sube, entonces responde 201 con `{id, filename, months, charges}` y `id` igual al `public_id`.
- Dada una cartola que dejó líneas ilegibles, cuando se sube, entonces la respuesta trae solo la cantidad de líneas malas, nunca su texto.
- Dada una sesión, cuando se sube, entonces la cartola queda a nombre de ese usuario.

## Fuera
- Nada.

## Notas
- Ninguna.
