# FOKKO-0084A-SUBS · Tokens de diseño en Tailwind

## Problema
El diseño usa tokens de color y tema claro/oscuro. El tema vive en el navegador (sin endpoint).

Esta idea cubre solo: Tokens de diseño en Tailwind. El resto del trabajo original está en las ideas hermanas `FOKKO-0084*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS.

## Alcance
- Tailwind y tokens

## Fuera
- La pantalla del selector → FOKKO-0110A-SUBS y FOKKO-0110B-SUBS.

## Features afectadas
- `frontend` (`lib/theme`, estilos)

## Dudas abiertas
- Los tokens salen del `tokens.json` del handoff Fokko Subs. Hace falta tener ese archivo a mano; no se inventan colores.

## Tareas
- [x] [TASK-001A Tokens en `:root`](TASK-001A-TOKENS-ROOT.md)
- [x] [TASK-001B Clases de Tailwind sobre los tokens](TASK-001B-TOKENS-CLASSES.md)
- [x] [TASK-001C Valores del tema oscuro](TASK-001C-TOKENS-DARK.md)
