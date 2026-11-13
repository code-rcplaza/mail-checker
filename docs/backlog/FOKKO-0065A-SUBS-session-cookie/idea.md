# FOKKO-0065A-SUBS · Cookie de sesión

## Problema
La sesión viaja en una cookie, así que hay que fijar sus flags y defenderse de CSRF. La app se sirve desde el mismo origen que la API, por lo que no necesita CORS.

Esta idea cubre solo: Cookie de sesión. El resto del trabajo original está en las ideas hermanas `FOKKO-0065*-SUBS`.

## Depende de
- FOKKO-0064-SUBS.

## Alcance
- Cookie de sesión

## Fuera
- Rate limit por IP: post-MVP.

## Features afectadas
- `internal/auth` (`adapters/http`)

## Dudas abiertas
- Se propone `SameSite=Lax` más una comprobación del encabezado `Origin` en `POST`, `PUT` y `DELETE`. Hace falta un error 403 para el origen no permitido (`common.forbidden`): ¿se agrega a `shared`?

## Tareas
- [ ] [TASK-001A Fijar los atributos de la cookie](TASK-001A-COOKIE-FLAGS.md)
- [ ] [TASK-001B Marcar `Secure` según el modo](TASK-001B-COOKIE-SECURE.md)
- [ ] [TASK-001C Borrar la cookie](TASK-001C-COOKIE-CLEAR.md)
