# FOKKO-0087-SUBS/TASK-001B · Bloquear el clic deshabilitado o cargando
Feature: `frontend/src/components/ui`



## Criterios
- Dado `disabled` o `loading`, cuando se hace clic, entonces no se llama al manejador.
- Dado `loading`, cuando se renderiza, entonces anuncia el estado de carga y mantiene su ancho.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
