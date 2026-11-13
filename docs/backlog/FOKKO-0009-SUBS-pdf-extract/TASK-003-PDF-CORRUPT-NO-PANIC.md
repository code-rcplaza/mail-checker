# FOKKO-0009-SUBS/TASK-003 · PDF dañado sin panic
Feature: `internal/statements/adapters/extract`

Como usuario quiero que un PDF dañado no tire abajo la app, para poder reintentar con otro archivo.

## Criterios
- Dado un PDF truncado o dañado, cuando se extrae, entonces devuelve un error con contexto y no hace panic.
- Dada una librería que hace panic con un PDF malformado, cuando se extrae, entonces el panic se convierte en error.

## Fuera
- Nada.

## Notas
- La librería puede hacer panic con entradas raras; el adaptador lo contiene.
