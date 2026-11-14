# FOKKO-0043-SUBS · Verificar el código (vigencia y 5 intentos)

## Problema
Verificar el código abre la sesión. Debe rechazar códigos vencidos, equivocados y agotados.

## Depende de
- FOKKO-0042-SUBS.

## Alcance
- Caso de uso `VerifyCode {email, code}`.

## Fuera
- Crear la sesión → FOKKO-0046-SUBS.

## Features afectadas
- `internal/auth` (`usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Código correcto](TASK-001-VERIFY-CODE-OK.md)
- [x] [TASK-002 Código equivocado](TASK-002-VERIFY-CODE-WRONG.md)
- [x] [TASK-003 Código vencido](TASK-003-VERIFY-CODE-EXPIRED.md)
