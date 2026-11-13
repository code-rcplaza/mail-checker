# FOKKO-0102-SUBS/TASK-002 · Confirmar con cadencia
Feature: `frontend/src/features/subscriptions/application`

Como usuario quiero decir "sí, es una suscripción" y elegir cada cuánto, para que la app lo recuerde.

## Criterios
- Dada una posible, cuando se toca "Sí", entonces exige elegir mensual o anual antes de enviar.
- Dada la cadencia elegida, cuando se confirma, entonces llama a `POST /api/subscriptions/decisions` y la tarjeta pasa a "agregada".
- Dado un error de la API, cuando falla, entonces la tarjeta vuelve a pendiente y muestra el mensaje.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
