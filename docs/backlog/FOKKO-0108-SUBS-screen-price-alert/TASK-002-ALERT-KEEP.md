# FOKKO-0108-SUBS/TASK-002 · Mantener
Feature: `frontend/src/features/subscriptions/application`

Como usuario quiero decir que acepto el nuevo precio, para que no me lo vuelva a mostrar.

## Criterios
- Dado "Mantener", cuando se toca, entonces llama a `POST /api/subscriptions/price-alerts/decision` y la alerta sale de pendientes.
- Dado un error de la API, cuando falla, entonces la alerta sigue visible con el mensaje.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
