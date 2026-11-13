# FOKKO-0065C-SUBS/TASK-001 · Sin cabeceras CORS
Feature: `internal/auth/adapters/http`

Como usuario quiero que mi API no se abra a otros sitios, para mi seguridad.

## Criterios
- Dada cualquier respuesta, cuando se inspeccionan sus cabeceras, entonces no hay `Access-Control-Allow-Origin`.
- Dada una petición `OPTIONS` de otro origen, cuando llega, entonces no se autoriza.

## Fuera
- Nada.

## Notas
- Ninguna.
