# FOKKO-0045-SUBS · Puerto de envío con doble de consola

## Problema
Enviar el código es un puerto (`CodeSender`). En desarrollo se usa un doble que escribe en consola; en producción, el adaptador real.

## Alcance
- Puerto `CodeSender` y adaptador de consola para desarrollo.

## Fuera
- Adaptador real → FOKKO-0048-SUBS.

## Features afectadas
- `internal/auth` (`ports`, `adapters/mailer`)

## Dudas abiertas
- El código en consola es un secreto: el adaptador de consola solo se usa en desarrollo y `main.go` se niega a usarlo en producción.

## Tareas
- [x] [TASK-001 Puerto de envío](TASK-001-CODE-SENDER-PORT.md)
- [x] [TASK-002 Envío por consola para desarrollo](TASK-002-CONSOLE-SENDER.md)
