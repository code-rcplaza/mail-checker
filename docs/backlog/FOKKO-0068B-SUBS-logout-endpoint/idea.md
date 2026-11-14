# FOKKO-0068B-SUBS · `POST /api/auth/logout`

## Problema
La app necesita saber quién está dentro (nombre, correo, iniciales del avatar) y cerrar sesión.

Esta idea cubre solo: `POST /api/auth/logout`. El resto del trabajo original está en las ideas hermanas `FOKKO-0068*-SUBS`.

## Depende de
- FOKKO-0047A-SUBS y FOKKO-0047B-SUBS, FOKKO-0050-SUBS y FOKKO-0064-SUBS.

## Alcance
- `POST /api/auth/logout`

## Fuera
- `DELETE /api/me` → FOKKO-0080-SUBS.

## Features afectadas
- `internal/auth` (`domain`, `usecases`, `adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 `POST /api/auth/logout`](TASK-001-LOGOUT-HANDLER.md)
