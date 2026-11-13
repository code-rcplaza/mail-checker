# FOKKO-0093-SUBS · `SideNav` de escritorio

## Problema
Desde 1024 px la app es una composición propia: barra lateral de 256 px y contenido en dos columnas. La tablet (768–1023 px) usa el layout de escritorio con barra inferior.

## Depende de
- FOKKO-0092-SUBS.
- FOKKO-0085C-SUBS.

## Alcance
- `SideNav`, el corte `lg:` del `AppShell` y la composición de dos columnas.

## Fuera
- Nada.

## Features afectadas
- `frontend` (`components/layout`)

## Dudas abiertas
- Las pruebas de jsdom no evalúan media queries: el corte se verifica con una revisión visual a 390, 768 y 1280 px, que se anota en el informe de cada pantalla.

## Tareas
- [ ] [TASK-001 `SideNav`](TASK-001-SIDE-NAV.md)
- [ ] [TASK-002 Corte de escritorio](TASK-002-SHELL-RESPONSIVE.md)
- [ ] [TASK-003 Contenido en dos columnas](TASK-003-DESKTOP-TWO-COLUMNS.md)
