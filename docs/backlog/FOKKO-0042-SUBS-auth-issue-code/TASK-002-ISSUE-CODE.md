# FOKKO-0042-SUBS/TASK-002 · Emitir un código
Feature: `internal/auth/usecases`

Como usuario quiero que pedir otro código anule el anterior, para no tener varios válidos.

## Criterios
- Dado un usuario, cuando se emite un código, entonces se guarda su hash y se devuelve el código en claro solo para el envío.
- Dado un usuario con un código vigente, cuando se emite otro, entonces el anterior deja de valer.
- Dado un error al guardar, cuando se emite, entonces devuelve el error con contexto y no se devuelve ningún código.

## Fuera
- Nada.

## Notas
- Ninguna.
