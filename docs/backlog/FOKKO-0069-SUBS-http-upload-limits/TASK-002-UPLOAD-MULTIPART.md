# FOKKO-0069-SUBS/TASK-002 · Formato de la subida
Feature: `internal/statements/adapters/http`

Como usuario quiero que se rechace una subida mal formada, para saber qué falló.

## Criterios
- Dada una petición que no es `multipart/form-data`, cuando se sube, entonces responde 400 con `common.invalid_request`.
- Dada una subida sin el campo `file`, cuando se sube, entonces responde 400 con `fields` indicando `file`.
- Dados dos archivos en el campo `file`, cuando se sube, entonces se rechaza.

## Fuera
- Nada.

## Notas
- Ninguna.
