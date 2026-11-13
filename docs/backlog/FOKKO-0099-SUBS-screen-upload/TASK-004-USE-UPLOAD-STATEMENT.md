# FOKKO-0099-SUBS/TASK-004 · `useUploadStatement`
Feature: `frontend/src/features/statements/application`

Como usuario quiero que la subida se haga sola al confirmar, para no repetir pasos.

## Criterios
- Dado un archivo, banco y fecha opcionales, cuando se confirma, entonces envía `multipart/form-data` con `file` y los campos presentes.
- Dada una respuesta 201, cuando termina, entonces entrega archivo, meses y cargos leídos.
- Dado un error de la API, cuando falla, entonces expone la llave y los parámetros.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
