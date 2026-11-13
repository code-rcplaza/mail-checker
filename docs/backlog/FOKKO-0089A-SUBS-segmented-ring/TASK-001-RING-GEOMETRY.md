# FOKKO-0089A-SUBS/TASK-001 · Geometría de segmentos
Feature: `frontend/src/components/ui`

Como usuario quiero que cada tajada del anillo sea proporcional, para entender mi gasto.

## Criterios
- Dados valores 50, 30 y 20, cuando se calcula, entonces los arcos miden 50 %, 30 % y 20 % del círculo menos la separación.
- Dado un solo segmento, cuando se calcula, entonces ocupa el círculo completo.
- Dado un total en cero, cuando se calcula, entonces no hay arcos y no hay división por cero.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
