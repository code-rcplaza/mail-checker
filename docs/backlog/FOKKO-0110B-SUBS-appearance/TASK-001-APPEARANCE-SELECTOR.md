# FOKKO-0110B-SUBS/TASK-001 · `AppearanceSelector`
Feature: `frontend/src/features/profile/presentation/components`

Como usuario quiero elegir claro, oscuro o sistema, para ver la app a mi gusto.

## Criterios
- Dadas las tres opciones, cuando se elige una, entonces cambia el tema al instante.
- Dada la elección, cuando se recarga, entonces se conserva.
- Dada la pantalla, cuando se renderiza, entonces no llama a ningún endpoint (vive en el navegador).

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
