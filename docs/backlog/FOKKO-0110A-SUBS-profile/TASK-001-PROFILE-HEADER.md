# FOKKO-0110A-SUBS/TASK-001 · `ProfileHeader`
Feature: `frontend/src/features/profile/presentation/components`

Como usuario quiero ver mi cuenta, para confirmar con cuál estoy.

## Criterios
- Dada una sesión, cuando se renderiza, entonces muestra avatar con iniciales, nombre y correo.
- Dado `GET /api/me` en carga, cuando se renderiza, entonces muestra un `Skeleton`.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
