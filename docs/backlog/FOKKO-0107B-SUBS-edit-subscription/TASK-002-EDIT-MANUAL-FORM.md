# FOKKO-0107B-SUBS/TASK-002 · Editar o quitar una manual
Feature: `frontend/src/features/subscriptions/application`

Como usuario quiero corregir o quitar una suscripción que agregué, para mantenerla al día.

## Criterios
- Dada una manual, cuando se edita y se guarda, entonces llama a `PUT /api/subscriptions/manual/{publicId}`.
- Dada una manual, cuando se quita, entonces pide confirmación y llama a `DELETE`.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
