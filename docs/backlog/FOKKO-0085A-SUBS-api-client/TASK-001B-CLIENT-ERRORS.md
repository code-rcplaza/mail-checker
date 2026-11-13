# FOKKO-0085A-SUBS/TASK-001B · Leer errores de la API
Feature: `frontend`



## Criterios
- Dado un cuerpo con `{error:{key,params,fields}}`, cuando se pide, entonces lanza un error tipado con esos tres datos.
- Dada una respuesta que no es JSON o está truncada, cuando se pide, entonces lanza un error genérico (`common.internal`).

## Fuera
- Nada.

## Notas
- Ninguna.
