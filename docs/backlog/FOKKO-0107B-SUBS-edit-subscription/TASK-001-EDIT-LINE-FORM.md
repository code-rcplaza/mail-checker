# FOKKO-0107B-SUBS/TASK-001 · Editar una línea detectada
Feature: `frontend/src/features/subscriptions/application`

Como usuario quiero ponerle nombre y corregir el monto o el día, para que mi resumen refleje lo real.

## Criterios
- Dada una línea, cuando se guarda un alias, entonces llama a `PUT /api/subscriptions/aliases`; alias vacío lo borra.
- Dado un monto o día distinto del detectado, cuando se guarda, entonces llama a `PUT /api/subscriptions/overrides`.
- Dada la edición desde Revisión, cuando se guarda, entonces vuelve a Revisión.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
