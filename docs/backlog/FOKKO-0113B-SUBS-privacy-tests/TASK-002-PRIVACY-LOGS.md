# FOKKO-0113B-SUBS/TASK-002 · Los logs no filtran datos
Feature: `cmd/api`

Como usuario quiero que mis datos no queden en los registros, para mi privacidad.

## Criterios
- Dado el log capturado durante una importación, cuando se revisa, entonces no contiene descriptores, montos ni nombres de archivo.
- Dado el log capturado durante un login, cuando se revisa, entonces no contiene el correo, el código ni la llave de Resend.
- Dado un error de importación, cuando se registra, entonces lleva solo la llave y números de línea.

## Fuera
- Nada.

## Notas
- Ninguna.
