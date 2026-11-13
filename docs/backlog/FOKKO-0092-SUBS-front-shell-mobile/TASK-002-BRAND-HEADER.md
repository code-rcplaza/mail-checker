# FOKKO-0092-SUBS/TASK-002 · `BrandHeader`
Feature: `frontend/src/components/layout`

Como usuario quiero ver la marca y mi avatar, para saber dónde estoy.

## Criterios
- Dada una sesión, cuando se renderiza, entonces muestra el wordmark, la pastilla del producto y el avatar con mis iniciales.
- Dada ninguna sesión, cuando se renderiza, entonces muestra el botón de tema en lugar del avatar.
- Dado el avatar, cuando se toca, entonces lleva a Perfil.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
