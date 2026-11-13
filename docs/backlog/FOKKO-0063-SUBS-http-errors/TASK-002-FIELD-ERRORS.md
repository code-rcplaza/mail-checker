# FOKKO-0063-SUBS/TASK-002 · Errores por campo
Feature: `internal/shared/httpserver`

Como usuario quiero ver qué campo está mal, para corregirlo.

## Criterios
- Dado un error de validación con campos, cuando se escribe, entonces el cuerpo es `{error:{key:"common.validation_failed", fields:{campo:llave}}}` con status 400.
- Dado un error sin campos, cuando se escribe, entonces el cuerpo no trae `fields`.

## Fuera
- Nada.

## Notas
- Ninguna.
