# FOKKO-0099-SUBS/TASK-001B · Rechazar archivos inválidos
Feature: `frontend/src/features/statements/presentation/components`



## Criterios
- Dado un archivo que no es PDF, cuando se elige, entonces se rechaza antes de enviar con el mensaje correspondiente.
- Dado un archivo de más de 10 MB, cuando se elige, entonces se rechaza antes de enviar.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
