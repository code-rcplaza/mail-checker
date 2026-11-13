# FOKKO-0064-SUBS/TASK-002 · Rutas sin sesión
Feature: `internal/auth/adapters/http`

Como usuario quiero poder entrar sin sesión previa, para abrir una.

## Criterios
- Dadas `/api/health` y `/api/auth/*`, cuando llegan sin cookie, entonces no se exige sesión.
- Dada cualquier otra ruta `/api/*`, cuando llega sin cookie, entonces responde 401.
- Dada una ruta que no existe, cuando llega sin cookie, entonces responde 404 en JSON (no filtra qué rutas existen detrás de la sesión).

## Fuera
- Nada.

## Notas
- Ninguna.
