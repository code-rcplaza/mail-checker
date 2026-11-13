# FOKKO-0100B-SUBS/TASK-001B · Mostrar los dos conteos del análisis
Feature: `frontend/src/features/statements/presentation/components`



## Criterios
- Dado `GET /api/summary` con 3 suscripciones seguras y 2 por confirmar, cuando se renderiza, entonces muestra 3 y 2.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
