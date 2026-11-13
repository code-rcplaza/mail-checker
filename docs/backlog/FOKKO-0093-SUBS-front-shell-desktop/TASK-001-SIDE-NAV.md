# FOKKO-0093-SUBS/TASK-001 · `SideNav`
Feature: `frontend/src/components/layout`

Como usuario quiero una barra lateral en pantalla grande, para navegar sin la barra inferior.

## Criterios
- Dada la barra lateral, cuando se renderiza, entonces trae marca, Resumen, Subs y el pie con avatar, nombre y correo.
- Dado el pie, cuando se activa, entonces lleva a Perfil.
- Dada la ruta actual, cuando se renderiza, entonces su opción queda marcada.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
