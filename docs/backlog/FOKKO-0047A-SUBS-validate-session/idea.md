# FOKKO-0047A-SUBS · Validar la sesión

## Problema
Cada ruta con sesión debe saber quién es el usuario; cerrar sesión debe invalidarla. Hoy no existe un error de "sin sesión".

Esta idea cubre solo: Validar la sesión. El resto del trabajo original está en las ideas hermanas `FOKKO-0047*-SUBS`.

## Depende de
- FOKKO-0046-SUBS.

## Alcance
- Error de falta de sesión
- Validar una sesión

## Fuera
- Middleware HTTP → FOKKO-0064-SUBS.

## Features afectadas
- `internal/shared` (`domain`)
- `internal/auth` (`usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Error de falta de sesión](TASK-001-UNAUTHORIZED-ERROR.md)
- [ ] [TASK-002 Validar una sesión](TASK-002-VALIDATE-SESSION.md)
