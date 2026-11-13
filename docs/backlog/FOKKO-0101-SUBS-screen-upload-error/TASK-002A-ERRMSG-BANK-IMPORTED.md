# FOKKO-0101-SUBS/TASK-002A · Mensajes de banco o cartola repetida
Feature: `frontend/src/features/statements/presentation/components`



## Criterios
- Dada `statements.bank_not_detected`, cuando se muestra, entonces sugiere elegir el banco a mano.
- Dada `statements.already_imported`, cuando se muestra, entonces dice que esa cartola ya está cargada.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
