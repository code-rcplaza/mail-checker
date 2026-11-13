# FOKKO-0115-SUBS/TASK-004 · `.dockerignore`
Feature: `.`

Como usuario quiero que la imagen no lleve secretos ni basura, para no filtrar datos en un repo público.

## Criterios
- Dado el contexto de build, entonces excluye `.git`, `.env`, `node_modules`, bases de datos locales y cartolas.
- Dado un `.env` local, cuando se construye la imagen, entonces su contenido no queda en ninguna capa.

## Fuera
- Nada.

## Notas
- Ninguna.
