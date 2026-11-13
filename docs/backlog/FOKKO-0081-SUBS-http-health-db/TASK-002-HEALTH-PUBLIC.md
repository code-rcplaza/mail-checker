# FOKKO-0081-SUBS/TASK-002 · La salud no pide sesión
Feature: `internal/shared/httpserver`

Como usuario quiero que el chequeo de la plataforma funcione sin cuenta, para que no falle por falta de sesión.

## Criterios
- Dada una petición sin cookie, cuando se pide `/api/health`, entonces no responde 401.
- Dado un `Origin` externo, cuando se pide `/api/health` con `GET`, entonces responde igual.

## Fuera
- Nada.

## Notas
- Ninguna.
