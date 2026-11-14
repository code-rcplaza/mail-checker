# FOKKO-0083C-SUBS · Corredor de pruebas del front

## Problema
No existe `frontend/`. Hace falta la base del proyecto con TypeScript estricto, lint y pruebas antes de escribir cualquier pantalla.

Esta idea cubre solo: Corredor de pruebas del front. El resto del trabajo original está en las ideas hermanas `FOKKO-0083*-SUBS`.

## Depende de
- FOKKO-0083B-SUBS.

## Alcance
- Pruebas de componentes

## Fuera
- Tailwind y tema → FOKKO-0084A-SUBS y FOKKO-0084B-SUBS.
- Servir el build desde Go → FOKKO-0114A-SUBS a FOKKO-0114D-SUBS.

## Features afectadas
- `frontend`

## Dudas abiertas
- Dependencias a aprobar: React, Vite, TypeScript, Biome y, para pruebas, Vitest con Testing Library. CLAUDE.md pide preguntar antes de agregar cada una.
- Comandos del front (`pnpm lint`, `pnpm typecheck`, `pnpm test`, `pnpm build`) se agregan a "Comandos" de `CLAUDE.md` con tu ok.

## Tareas
- [x] [TASK-001 Pruebas de componentes](TASK-001-FRONT-TEST-RUNNER.md)
