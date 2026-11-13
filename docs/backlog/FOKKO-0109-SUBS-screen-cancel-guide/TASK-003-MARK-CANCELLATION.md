# FOKKO-0109-SUBS/TASK-003 · Marcar "por cancelar"
Feature: `frontend/src/features/subscriptions/application`

Como usuario quiero marcar la suscripción que voy a cancelar, para recordarlo.

## Criterios
- Dado "Ya la cancelé" o "La voy a cancelar", cuando se toca, entonces llama a `POST /api/subscriptions/cancellations`.
- Dada la marca guardada, cuando se renderiza, entonces la línea muestra el estado `CancellationStatus` y el botón queda deshabilitado.
- Dado un error de la API, cuando falla, entonces no cambia el estado y muestra el mensaje.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
