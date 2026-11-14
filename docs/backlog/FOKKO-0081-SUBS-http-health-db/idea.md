# FOKKO-0081-SUBS · `GET /api/health` también revisa la base (para el contenedor)

## Problema
El contenedor necesita saber si la app está sana, y eso incluye poder hablar con la base. `/api/health` ya existe y hoy solo responde "ok".

## Alcance
- `GET /api/health` comprueba la base con un timeout corto y responde 503 si no puede.

## Fuera
- Un `/healthz` aparte: no hace falta (YAGNI).

## Features afectadas
- `internal/shared` (`httpserver`, `sqlite`)

## Dudas abiertas
- Hace falta un status 503 mapeado (`common.unavailable`).

## Tareas
- [x] [TASK-001 La salud comprueba la base](TASK-001-HEALTH-DB-PING.md)
- [x] [TASK-002 La salud no pide sesión](TASK-002-HEALTH-PUBLIC.md)
