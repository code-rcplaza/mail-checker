# FOKKO-0061-SUBS/TASK-002 · Caso de uso de limpieza
Feature: `internal/auth/usecases`

Como usuario quiero que la limpieza use el reloj de la app, para que sea predecible.

## Criterios
- Dado el reloj en T, cuando se ejecuta, entonces pide borrar lo anterior a T.
- Dado un error del puerto, cuando se ejecuta, entonces devuelve el error con contexto.

## Fuera
- Nada.

## Notas
- Ninguna.
