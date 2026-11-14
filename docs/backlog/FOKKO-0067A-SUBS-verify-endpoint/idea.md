# FOKKO-0067A-SUBS · `POST /api/auth/code/verify`

## Problema
Verificar el código abre la sesión (cookie); reenviar pide otro código.

Esta idea cubre solo: `POST /api/auth/code/verify`. El resto del trabajo original está en las ideas hermanas `FOKKO-0067*-SUBS`.

## Depende de
- FOKKO-0043-SUBS, FOKKO-0046-SUBS, FOKKO-0065A-SUBS a FOKKO-0065C-SUBS y FOKKO-0066A-SUBS y FOKKO-0066B-SUBS.

## Alcance
- `POST /api/auth/code/verify`
- Formato del código en el borde

## Fuera
- Nada.

## Features afectadas
- `internal/auth` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 `POST /api/auth/code/verify`](TASK-001-VERIFY-HANDLER.md)
- [x] [TASK-002 Formato del código en el borde](TASK-002-CODE-FORMAT-AT-EDGE.md)
