# FOKKO-0065C-SUBS · Sin cabeceras CORS

## Problema
La sesión viaja en una cookie, así que hay que fijar sus flags y defenderse de CSRF. La app se sirve desde el mismo origen que la API, por lo que no necesita CORS.

Esta idea cubre solo: Sin cabeceras CORS. El resto del trabajo original está en las ideas hermanas `FOKKO-0065*-SUBS`.

## Depende de
- FOKKO-0064-SUBS.

## Alcance
- Sin cabeceras CORS

## Fuera
- Rate limit por IP: post-MVP.

## Features afectadas
- `internal/auth` (`adapters/http`)

## Dudas abiertas
- Se propone `SameSite=Lax` más una comprobación del encabezado `Origin` en `POST`, `PUT` y `DELETE`. Hace falta un error 403 para el origen no permitido (`common.forbidden`): ¿se agrega a `shared`?

## Tareas
- [x] [TASK-001 Sin cabeceras CORS](TASK-001-NO-CORS.md)
