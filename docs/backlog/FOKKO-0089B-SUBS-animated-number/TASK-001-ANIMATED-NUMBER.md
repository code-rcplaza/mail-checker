# FOKKO-0089B-SUBS/TASK-001 · `AnimatedNumber`
Feature: `frontend/src/components/ui`

Como usuario quiero que el total aparezca con vida, para notar el resultado.

## Criterios
- Dado un valor, cuando se monta, entonces cuenta de 0 al valor en 0,9 s y termina exactamente en el valor.
- Dada la preferencia `prefers-reduced-motion`, cuando se monta, entonces muestra el valor final sin animar.
- Dado un valor nuevo, cuando cambia, entonces anima desde el valor actual al nuevo.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
