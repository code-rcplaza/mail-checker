# FOKKO-0069-SUBS/TASK-003 · Tiempo máximo
Feature: `internal/statements/adapters/http`

Como usuario quiero que una subida colgada no bloquee a los demás, para que la app siga disponible.

## Criterios
- Dado un cliente que envía muy lento, cuando pasa el tiempo máximo, entonces se corta la petición.
- Dada una lectura del PDF que supera el tiempo máximo, cuando se procesa, entonces se cancela por contexto y responde error.

## Fuera
- Nada.

## Notas
- Ninguna.
