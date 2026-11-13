# FOKKO-0048-SUBS/TASK-001 · Enviar el código por Resend
Feature: `internal/auth/adapters/mailer`

Como usuario quiero recibir mi código en mi bandeja, para poder entrar.

## Criterios
- Dado un servidor de prueba (`httptest`), cuando se envía, entonces hace `POST` a `/emails` con `Authorization: Bearer <llave>` y JSON con remitente, destinatario, asunto y texto que incluye el código.
- Dado un envío correcto (2xx), cuando se envía, entonces devuelve sin error.

## Fuera
- Nada.

## Notas
- Ninguna.
