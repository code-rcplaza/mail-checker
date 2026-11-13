# FOKKO-0107A-SUBS/TASK-002 · Agregar una manual
Feature: `frontend/src/features/subscriptions/application`

Como usuario quiero agregar una suscripción que no aparece en mi cartola, para ver todo mi gasto.

## Criterios
- Dados datos válidos, cuando se envía, entonces llama a `POST /api/subscriptions/manual` y vuelve a Subs con la nueva fila.
- Dado un 400 con `fields`, cuando se envía, entonces cada error aparece en su campo.
- Dada una pantalla móvil y una de escritorio, cuando se revisa, entonces ambas siguen el diseño.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
