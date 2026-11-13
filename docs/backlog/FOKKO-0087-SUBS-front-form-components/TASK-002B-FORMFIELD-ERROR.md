# FOKKO-0087-SUBS/TASK-002B · Errores de `FormField`
Feature: `frontend/src/components/ui`



## Criterios
- Dado un error, cuando se renderiza, entonces el campo queda marcado como inválido y el error se anuncia.
- Dada una ayuda y un error a la vez, cuando se renderiza, entonces se muestra el error primero.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
