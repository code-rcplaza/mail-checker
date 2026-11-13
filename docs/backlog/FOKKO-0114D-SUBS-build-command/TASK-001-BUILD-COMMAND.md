# FOKKO-0114D-SUBS/TASK-001 · Un solo comando de build
Feature: `backend`

Como usuario quiero construir todo con un comando, para no olvidar un paso.

## Criterios
- Dado el comando de build, cuando se corre en un checkout limpio, entonces construye el front y luego el binario de Go con el front adentro.
- Dado un fallo del build del front, cuando ocurre, entonces el comando termina con error y no produce un binario a medias.
- Dado el comando, entonces queda en "Comandos" de `CLAUDE.md`.

## Fuera
- Nada.

## Notas
- Ninguna.
