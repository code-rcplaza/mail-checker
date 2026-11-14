# FOKKO-0066B-SUBS · `POST /api/auth/login`

## Problema
Registro y login piden un código al correo. El borde HTTP valida y rechaza lo inválido antes de que llegue adentro.

Esta idea cubre solo: `POST /api/auth/login`. El resto del trabajo original está en las ideas hermanas `FOKKO-0066*-SUBS`.

## Depende de
- FOKKO-0041-SUBS, FOKKO-0064-SUBS y FOKKO-0063-SUBS.

## Alcance
- `POST /api/auth/login`

## Fuera
- Verificar y reenviar → FOKKO-0067A-SUBS y FOKKO-0067B-SUBS.

## Features afectadas
- `internal/auth` (`usecases`, `adapters/http`)

## Dudas abiertas
- Respuesta de éxito: se propone `204 No Content`, porque el resultado es el correo que llega.

## Tareas
- [x] [TASK-001A Responder 204 al entrar](TASK-001A-LOGIN-204.md)
- [x] [TASK-001B Rechazar correo inválido o desconocido](TASK-001B-LOGIN-REJECT.md)
- [x] [TASK-001C Responder 429 por límite](TASK-001C-LOGIN-429.md)
