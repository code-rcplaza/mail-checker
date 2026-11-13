# FOKKO-0048-SUBS/TASK-003 · Configuración del envío
Feature: `cmd/api`

Como usuario quiero que la app no arranque a medias, para no descubrir en el primer login que falta la llave.

## Criterios
- Dado el modo producción sin `RESEND_API_KEY` o sin remitente, cuando se carga la configuración, entonces el arranque falla con un mensaje claro.
- Dado el modo desarrollo sin llave, cuando se carga la configuración, entonces usa el envío por consola.
- Dada la llave cargada, cuando se imprime o se loguea la configuración, entonces la llave no aparece.

## Fuera
- Nada.

## Notas
- Ninguna.
