# FOKKO-0106A-SUBS/TASK-002 · `PriceChangeBlock`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver cuánto subió, para decidir si sigo.

## Criterios
- Dada una alza, cuando se renderiza, entonces muestra antes, después, desde cuándo y el impacto mensual y anual.
- Dada una alza marcada como cambio de plan, cuando se renderiza, entonces lo dice.
- Dada una línea sin alza, cuando se renderiza, entonces el bloque no aparece (caso sin alza).

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
