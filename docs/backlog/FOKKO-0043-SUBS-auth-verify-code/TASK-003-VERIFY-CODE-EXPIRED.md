# FOKKO-0043-SUBS/TASK-003 · Código vencido
Feature: `internal/auth/usecases`

Como usuario quiero un mensaje claro si mi código venció, para pedir otro.

## Criterios
- Dado un código con más de 10 minutos, cuando se verifica, entonces devuelve `ErrInvalid` con `auth.code_expired`.
- Dado un código a exactamente 10 minutos, cuando se verifica, entonces se documenta y prueba de qué lado cae.
- Dado un correo sin código emitido, cuando se verifica, entonces devuelve `auth.code_invalid` (no revela si el correo existe).

## Fuera
- Nada.

## Notas
- Ninguna.
