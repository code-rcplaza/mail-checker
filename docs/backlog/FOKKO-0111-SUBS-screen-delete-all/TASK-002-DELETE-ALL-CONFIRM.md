# FOKKO-0111-SUBS/TASK-002 · Confirmación
Feature: `frontend/src/features/profile/presentation/components`

Como usuario quiero una confirmación clara, para no borrar por accidente.

## Criterios
- Dado "Eliminar todo", cuando se toca, entonces muestra lo que se borra y pide confirmar.
- Dado el diálogo, cuando se presiona Esc o "Cancelar", entonces se cierra sin llamar a la API.
- Dado el diálogo abierto, entonces el foco queda dentro y el scroll de fondo bloqueado.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
