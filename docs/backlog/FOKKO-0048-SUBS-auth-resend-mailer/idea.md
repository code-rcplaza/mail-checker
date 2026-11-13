# FOKKO-0048-SUBS · Envío del código con Resend (`net/http`, timeout, llave por entorno)

## Problema
El código llega por correo con Resend. Se llama por HTTP con `net/http`, sin SDK, con timeout, y la llave viene por variable de entorno.

## Depende de
- FOKKO-0045-SUBS.

## Alcance
- Adaptador `CodeSender` con la API de Resend.
- Configuración (`RESEND_API_KEY`, remitente) que detiene el arranque si falta.

## Fuera
- Plantillas de correo bonitas: texto simple.

## Features afectadas
- `internal/auth` (`adapters/mailer`)
- `cmd/api` (`config.go`)

## Dudas abiertas
- Hay que tener un dominio verificado en Resend y definir el remitente (`MAIL_FROM`); sin eso solo llegan correos a la cuenta dueña.
- Texto del correo en español: asunto y cuerpo se definen al escribir la tarea.

## Tareas
- [ ] [TASK-001 Enviar el código por Resend](TASK-001-RESEND-SEND.md)
- [ ] [TASK-002 Errores del envío](TASK-002-RESEND-ERRORS.md)
- [ ] [TASK-003 Configuración del envío](TASK-003-RESEND-CONFIG.md)
