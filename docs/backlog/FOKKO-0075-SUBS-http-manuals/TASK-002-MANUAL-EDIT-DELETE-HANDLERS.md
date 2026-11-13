# FOKKO-0075-SUBS/TASK-002 · `PUT`/`DELETE` de una manual
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero corregir o quitar una manual, para mantenerla al día.

## Criterios
- Dado un `publicId` propio y datos válidos, cuando se edita, entonces responde 204.
- Dado un `publicId` propio, cuando se borra, entonces responde 204.
- Dado un `publicId` ajeno, inexistente o con formato inválido, cuando se edita o se borra, entonces responde 404.

## Fuera
- Nada.

## Notas
- Ninguna.
