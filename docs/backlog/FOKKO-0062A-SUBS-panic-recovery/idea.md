# FOKKO-0062A-SUBS · Recuperar panics

## Problema
Antes de exponer rutas con datos de usuarios, la base HTTP debe ser robusta: sin caídas por panics, con logs útiles que no filtren datos y con timeouts. Ya existe un router con `/api/health` y 404 en JSON (idea FOKKO-0001-SUBS); esta idea cubre lo que falta.

Esta idea cubre solo: Recuperar panics. El resto del trabajo original está en las ideas hermanas `FOKKO-0062*-SUBS`.

## Alcance
- Panic convertido en 500

## Fuera
- Sesión, cookies y CSRF → FOKKO-0064-SUBS y FOKKO-0065A-SUBS a FOKKO-0065C-SUBS.

## Features afectadas
- `internal/shared` (`httpserver`)
- `cmd/api`

## Dudas abiertas
- Antes de escribir las tareas se revisa qué de esto ya existe en `router.go` y `run.go`, para no duplicarlo.

## Tareas
- [ ] [TASK-001 Panic convertido en 500](TASK-001-PANIC-RECOVERY.md)
