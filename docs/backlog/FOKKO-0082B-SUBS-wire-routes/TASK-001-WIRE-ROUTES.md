# FOKKO-0082B-SUBS/TASK-001 · Registrar las rutas
Feature: `cmd/api`

Como usuario quiero que cada ruta del MVP exista y esté protegida, para usar toda la app.

## Criterios
- Dada la lista de rutas de `ARCHITECTURE.md`, cuando se arma el router, entonces cada una responde (no 404).
- Dadas todas las rutas salvo `/api/health` y `/api/auth/*`, cuando se piden sin sesión, entonces responden 401.
- Dada una ruta post-MVP (`/api/accounts`, `/api/sync`), cuando se pide, entonces no existe.

## Fuera
- Nada.

## Notas
- Ninguna.
