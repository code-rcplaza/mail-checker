# FOKKO-0067A-SUBS/TASK-002 · Formato del código en el borde
Feature: `internal/auth/adapters/http`

Como usuario quiero un aviso inmediato si escribí mal el código, para corregirlo sin gastar un intento.

## Criterios
- Dado un código que no son 6 dígitos, cuando se envía, entonces responde 400 con `fields` y no cuenta como intento.
- Dado un código con espacios alrededor, cuando se envía, entonces se recorta antes de validar.

## Fuera
- Nada.

## Notas
- Ninguna.
