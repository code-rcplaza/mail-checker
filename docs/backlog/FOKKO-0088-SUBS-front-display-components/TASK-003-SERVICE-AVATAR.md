# FOKKO-0088-SUBS/TASK-003 · `ServiceAvatar`
Feature: `frontend/src/components/ui`

Como usuario quiero reconocer cada servicio de un vistazo, para ubicarme.

## Criterios
- Dado el nombre "Netflix", cuando se renderiza, entonces muestra "N".
- Dado el mismo nombre dos veces, cuando se renderiza, entonces el color es el mismo.
- Dado un nombre vacío, cuando se renderiza, entonces muestra un símbolo neutro y no falla.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
