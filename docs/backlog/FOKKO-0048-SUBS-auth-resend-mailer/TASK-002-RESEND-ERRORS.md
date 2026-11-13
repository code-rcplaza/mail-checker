# FOKKO-0048-SUBS/TASK-002 · Errores del envío
Feature: `internal/auth/adapters/mailer`

Como usuario quiero que un fallo del proveedor no me deje colgado, para poder reintentar.

## Criterios
- Dada una respuesta no 2xx, cuando se envía, entonces devuelve un error con el status y sin el cuerpo, el código ni el correo.
- Dado un servidor que no responde, cuando se envía, entonces corta por timeout.
- Dado un contexto cancelado, cuando se envía, entonces devuelve el error sin enviar.

## Fuera
- Nada.

## Notas
- Ninguna.
