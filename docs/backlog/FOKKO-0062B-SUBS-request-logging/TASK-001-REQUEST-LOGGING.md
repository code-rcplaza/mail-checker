# FOKKO-0062B-SUBS/TASK-001 · Log de peticiones sin datos personales
Feature: `internal/shared/httpserver`

Como usuario quiero que mis datos no queden en los logs, para mi privacidad.

## Criterios
- Dada una petición, cuando termina, entonces se loguea método, ruta, status y duración.
- Dada una petición con query (`?lineKey=…`) o cuerpo con correo, cuando se loguea, entonces ninguno de los dos aparece.
- Dada una respuesta 5xx, cuando se loguea, entonces sube el nivel del log.

## Fuera
- Nada.

## Notas
- Ninguna.
