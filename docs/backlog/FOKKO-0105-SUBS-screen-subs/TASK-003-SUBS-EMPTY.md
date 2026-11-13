# FOKKO-0105-SUBS/TASK-003 · Subs vacío
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero una guía si no tengo suscripciones, para saber cómo seguir.

## Criterios
- Dado ningún dato, cuando se renderiza, entonces muestra el estado vacío con las acciones "Subir cartola" y "Agregar".

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
