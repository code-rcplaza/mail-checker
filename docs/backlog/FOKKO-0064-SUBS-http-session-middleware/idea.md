# FOKKO-0064-SUBS · Middleware de sesión

## Problema
Toda ruta exige sesión salvo `/api/health` y `/api/auth/*`. El `userID` sale de la sesión en el borde HTTP.

## Depende de
- FOKKO-0047A-SUBS y FOKKO-0047B-SUBS y FOKKO-0063-SUBS.

## Alcance
- Middleware que lee la cookie, valida la sesión y deja el `userID` para los handlers.

## Fuera
- Cómo se fija la cookie → FOKKO-0065A-SUBS a FOKKO-0065C-SUBS.

## Features afectadas
- `internal/auth` (`adapters/http`)
- `internal/shared` (`httpserver`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Exigir sesión](TASK-001-SESSION-MIDDLEWARE.md)
- [ ] [TASK-002 Rutas sin sesión](TASK-002-PUBLIC-ROUTES.md)
- [ ] [TASK-003 El `userID` en el contexto](TASK-003-USER-ID-IN-CONTEXT.md)
