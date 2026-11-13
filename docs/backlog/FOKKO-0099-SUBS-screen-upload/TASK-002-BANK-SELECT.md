# FOKKO-0099-SUBS/TASK-002 · `BankSelect`
Feature: `frontend/src/features/statements/presentation/components`

Como usuario quiero indicar mi banco si la app no lo detecta, para que lea mi cartola.

## Criterios
- Dada la lista de `GET /api/banks`, cuando se renderiza, entonces ofrece cada banco con su nombre traducido y la opción "Detectar automáticamente".
- Dada la opción automática, cuando se envía, entonces no se manda el campo `bank`.
- Dado un fallo al cargar la lista, cuando ocurre, entonces el selector queda en automático y no bloquea la subida.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
