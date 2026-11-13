# FOKKO-0090-SUBS/TASK-001 · `Toast`
Feature: `frontend/src/components/ui`

Como usuario quiero un aviso breve tras una acción, para saber que se hizo.

## Criterios
- Dado un aviso, cuando se muestra, entonces se anuncia a lectores de pantalla y se cierra solo.
- Dado un aviso, cuando se cierra con el botón o con Esc, entonces desaparece.
- Dados dos avisos seguidos, cuando se muestran, entonces el nuevo reemplaza al anterior sin apilarse sin límite.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
