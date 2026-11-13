# FOKKO-0049-SUBS/TASK-002 · Comprobar que lo generado está al día
Feature: `backend`

Como usuario quiero detectar si alguien cambió el SQL y olvidó regenerar, para no desplegar código viejo.

## Criterios
- Dado el SQL sin cambios, cuando se corre la comprobación (`sqlc diff` o equivalente), entonces no hay diferencias.
- Dado un cambio en un `queries.sql` sin regenerar, cuando se corre la comprobación, entonces falla.
- Dado el flujo de verificación del proyecto, entonces la comprobación queda en "Comandos" de `CLAUDE.md`.

## Fuera
- Nada.

## Notas
- Ninguna.
