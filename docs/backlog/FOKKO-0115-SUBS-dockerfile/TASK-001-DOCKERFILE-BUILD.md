# FOKKO-0115-SUBS/TASK-001 · Dockerfile multi-stage
Feature: `.`

Como usuario quiero una imagen que contenga toda la app, para desplegarla en cualquier lado.

## Criterios
- Dado el Dockerfile, cuando se construye, entonces compila el front y el back (sin cgo) y la imagen final solo contiene el binario.
- Dada la imagen, cuando se inspecciona, entonces no contiene código fuente, el `.git` ni archivos `.env`.
- Dada la imagen, cuando arranca con la configuración mínima, entonces `GET /api/health` responde 200.

## Fuera
- Nada.

## Notas
- Ninguna.
