# FOKKO-0092-SUBS/TASK-003 · `BottomNav`
Feature: `frontend/src/components/layout`

Como usuario quiero cambiar entre Resumen y Subs con un toque, para moverme rápido.

## Criterios
- Dadas las rutas Resumen y Subs, cuando se renderiza, entonces la actual queda marcada.
- Dado un toque en la otra opción, cuando se hace, entonces navega a ella.
- Dado cada elemento, entonces mide al menos 44 px y tiene nombre accesible.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
