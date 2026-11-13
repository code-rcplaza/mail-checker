# FOKKO-0092-SUBS/TASK-001 · `AppShell` móvil
Feature: `frontend/src/components/layout`

Como usuario quiero una estructura estable en cada pantalla, para orientarme.

## Criterios
- Dado el shell, cuando se renderiza, entonces tiene la zona superior, la hoja con esquinas de 28 px superpuesta y el contenido.
- Dado contenido más largo que la pantalla, cuando se desplaza, entonces solo se mueve la hoja y la navegación inferior queda fija.
- Dado un navegador con zona segura (notch), cuando se renderiza, entonces respeta las áreas seguras.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
