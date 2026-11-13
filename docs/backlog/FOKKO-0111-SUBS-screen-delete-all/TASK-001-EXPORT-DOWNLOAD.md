# FOKKO-0111-SUBS/TASK-001 · Descargar la copia
Feature: `frontend/src/features/profile/application`

Como usuario quiero llevarme mis datos antes de borrar, para no perderlos.

## Criterios
- Dada la pantalla, cuando se renderiza, entonces ofrece "Descargar copia primero" con CSV por defecto y JSON como alternativa.
- Dada la descarga, cuando se toca, entonces pide `GET /api/export?format=…` y el navegador guarda el archivo.
- Dado un error de la API, cuando falla, entonces avisa y no bloquea la eliminación.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
