# FOKKO-0083A-SUBS · Proyecto Vite con React

## Problema
No existe `frontend/`. Hace falta la base del proyecto con TypeScript estricto, lint y pruebas antes de escribir cualquier pantalla.

Esta idea cubre solo: Proyecto Vite con React. El resto del trabajo original está en las ideas hermanas `FOKKO-0083*-SUBS`.

## Alcance
- Proyecto Vite y React

## Fuera
- Tailwind y tema → FOKKO-0084A-SUBS y FOKKO-0084B-SUBS.
- Servir el build desde Go → FOKKO-0114A-SUBS a FOKKO-0114D-SUBS.

## Features afectadas
- `frontend`

## Dudas abiertas
- Dependencias a aprobar: React, Vite, TypeScript, Biome y, para pruebas, Vitest con Testing Library. CLAUDE.md pide preguntar antes de agregar cada una.
- Comandos del front (`pnpm lint`, `pnpm typecheck`, `pnpm test`, `pnpm build`) se agregan a "Comandos" de `CLAUDE.md` con tu ok.

## Tareas
- [x] [TASK-001A Compilar el sitio](TASK-001A-SCAFFOLD-BUILD.md)
- [x] [TASK-001B Proxy de desarrollo hacia la API](TASK-001B-SCAFFOLD-PROXY.md)
- [x] [TASK-001C Montar la app](TASK-001C-SCAFFOLD-MOUNT.md)
