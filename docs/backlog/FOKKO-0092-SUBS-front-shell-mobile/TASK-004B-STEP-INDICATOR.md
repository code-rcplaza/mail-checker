# FOKKO-0092-SUBS/TASK-004B · `StepIndicator`
Feature: `frontend/src/components/layout`



## Criterios
- Dado `StepIndicator` con paso 2 de 3 y nombre "Análisis", cuando se renderiza, entonces dice "Paso 2 de 3 · Análisis".
- Dado un paso fuera de rango, cuando se renderiza, entonces no se rompe.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
