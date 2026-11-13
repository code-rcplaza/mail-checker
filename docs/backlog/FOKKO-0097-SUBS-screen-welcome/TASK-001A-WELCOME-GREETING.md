# FOKKO-0097-SUBS/TASK-001A · Saludar con el nombre
Feature: `frontend/src/features/auth/presentation/components`



## Criterios
- Dada una sesión con nombre "Ana Pérez", cuando se renderiza, entonces saluda con el nombre.
- Dada la pantalla, cuando se renderiza, entonces no hace ninguna llamada a la API distinta de `GET /api/me`.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
