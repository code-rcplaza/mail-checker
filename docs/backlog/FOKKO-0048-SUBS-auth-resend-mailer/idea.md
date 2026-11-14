# FOKKO-0048-SUBS · Envío del código con Resend (`net/http`, timeout, llave por entorno)

## Problema
El código llega por correo con Resend. Se llama por HTTP con `net/http`, sin SDK, con timeout, y la llave viene por variable de entorno.

## Depende de
- FOKKO-0045-SUBS.

## Alcance
- Adaptador `CodeSender` con la API de Resend.
- Configuración (`MAILCHECKER_ENV`, `MAILCHECKER_RESEND_API_KEY`, `MAILCHECKER_MAIL_FROM`) que detiene el arranque si falta.

## Fuera
- Plantillas de correo bonitas: texto simple.

## Features afectadas
- `internal/auth` (`adapters/mailer`)
- `cmd/api` (`config.go`)

## Dudas abiertas
- Hay que tener un dominio verificado en Resend y definir el remitente (`MAILCHECKER_MAIL_FROM`); sin eso solo llegan correos a la cuenta dueña. Lo configura el usuario antes de probar con correos reales.
- Texto del correo: decidido. Asunto "Tu código de acceso a Fokko"; cuerpo "Tu código es 123456. Vence en 10 minutos. Si no lo pediste, ignora este correo."

## Tareas
- [x] [TASK-001 Enviar el código por Resend](TASK-001-RESEND-SEND.md)
- [x] [TASK-002 Errores del envío](TASK-002-RESEND-ERRORS.md)
- [x] [TASK-003 Configuración del envío](TASK-003-RESEND-CONFIG.md)
