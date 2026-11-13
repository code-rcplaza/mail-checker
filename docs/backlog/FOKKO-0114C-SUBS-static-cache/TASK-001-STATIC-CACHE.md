# FOKKO-0114C-SUBS/TASK-001 · Cabeceras de caché
Feature: `internal/shared/httpserver`

Como usuario quiero que la app cargue rápido y se actualice sin trucos, para no ver versiones viejas.

## Criterios
- Dado un archivo con hash, cuando se sirve, entonces lleva caché larga e inmutable.
- Dado `index.html`, cuando se sirve, entonces lleva `no-cache`.
- Dadas las respuestas de la API, cuando se sirven, entonces no llevan caché pública.

## Fuera
- Nada.

## Notas
- Ninguna.
