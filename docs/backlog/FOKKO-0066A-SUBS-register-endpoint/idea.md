# FOKKO-0066A-SUBS · `POST /api/auth/register`

## Problema
Registro y login piden un código al correo. El borde HTTP valida y rechaza lo inválido antes de que llegue adentro.

Esta idea cubre solo: `POST /api/auth/register`. El resto del trabajo original está en las ideas hermanas `FOKKO-0066*-SUBS`.

## Depende de
- FOKKO-0041-SUBS, FOKKO-0064-SUBS y FOKKO-0063-SUBS.

## Alcance
- Registrarse y recibir el código
- `POST /api/auth/register`

## Fuera
- Verificar y reenviar → FOKKO-0067A-SUBS y FOKKO-0067B-SUBS.

## Features afectadas
- `internal/auth` (`usecases`, `adapters/http`)

## Dudas abiertas
- Respuesta de éxito: se propone `204 No Content`, porque el resultado es el correo que llega.

## Tareas
- [x] [TASK-001A Registrar enviando el código](TASK-001A-REGISTER-SENDS-CODE.md)
- [x] [TASK-001B Registrar un correo ya registrado](TASK-001B-REGISTER-CONFLICT.md)
- [x] [TASK-001C Propagar un fallo del envío](TASK-001C-REGISTER-SEND-FAILS.md)
- [x] [TASK-002A Responder 204 al registrar](TASK-002A-REGISTER-204.md)
- [x] [TASK-002B Responder 400 por datos inválidos](TASK-002B-REGISTER-400.md)
- [x] [TASK-002C Responder 204 por correo repetido](TASK-002C-REGISTER-409.md)
