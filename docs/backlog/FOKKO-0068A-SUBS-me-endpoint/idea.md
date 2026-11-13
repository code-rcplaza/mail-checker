# FOKKO-0068A-SUBS · `GET /api/me`

## Problema
La app necesita saber quién está dentro (nombre, correo, iniciales del avatar) y cerrar sesión.

Esta idea cubre solo: `GET /api/me`. El resto del trabajo original está en las ideas hermanas `FOKKO-0068*-SUBS`.

## Depende de
- FOKKO-0047A-SUBS y FOKKO-0047B-SUBS, FOKKO-0050-SUBS y FOKKO-0064-SUBS.

## Alcance
- Iniciales del avatar
- `GET /api/me`

## Fuera
- `DELETE /api/me` → FOKKO-0080-SUBS.

## Features afectadas
- `internal/auth` (`domain`, `usecases`, `adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Iniciales de uno o dos nombres](TASK-001A-INITIALS-BASIC.md)
- [ ] [TASK-001B Iniciales con acentos o emoji](TASK-001B-INITIALS-UNICODE.md)
- [ ] [TASK-002 `GET /api/me`](TASK-002-GET-ME.md)
