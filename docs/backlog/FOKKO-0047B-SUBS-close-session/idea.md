# FOKKO-0047B-SUBS · Cerrar la sesión

## Problema
Cada ruta con sesión debe saber quién es el usuario; cerrar sesión debe invalidarla. Hoy no existe un error de "sin sesión".

Esta idea cubre solo: Cerrar la sesión. El resto del trabajo original está en las ideas hermanas `FOKKO-0047*-SUBS`.

## Depende de
- FOKKO-0046-SUBS.

## Alcance
- Cerrar la sesión

## Fuera
- Middleware HTTP → FOKKO-0064-SUBS.

## Features afectadas
- `internal/shared` (`domain`)
- `internal/auth` (`usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Cerrar la sesión](TASK-001-CLOSE-SESSION.md)
