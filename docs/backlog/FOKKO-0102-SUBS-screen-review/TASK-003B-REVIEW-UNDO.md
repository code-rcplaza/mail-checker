# FOKKO-0102-SUBS/TASK-003B · Deshacer una decisión
Feature: `frontend/src/features/subscriptions/application`



## Criterios
- Dada una tarjeta resuelta, cuando se toca "Deshacer", entonces llama a `DELETE /api/subscriptions/decisions` y vuelve a pendiente.
- Dado el deshacer en curso, cuando se toca otra vez, entonces no se envían dos llamadas.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
