# FOKKO-0079-SUBS/TASK-004 · `GET /api/export`
Feature: `internal/auth/adapters/http`

Como usuario quiero descargar mi copia desde Perfil, para guardarla antes de eliminar todo.

## Criterios
- Dado `?format=csv` o sin formato, cuando se pide, entonces responde CSV con `Content-Disposition: attachment`.
- Dado `?format=json`, cuando se pide, entonces responde JSON descargable.
- Dado otro formato, cuando se pide, entonces responde 400 con `fields`.

## Fuera
- Nada.

## Notas
- Ninguna.
