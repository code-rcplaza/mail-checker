# FOKKO-0062B-SUBS · Log de peticiones sin datos personales

## Problema
Antes de exponer rutas con datos de usuarios, la base HTTP debe ser robusta: sin caídas por panics, con logs útiles que no filtren datos y con timeouts. Ya existe un router con `/api/health` y 404 en JSON (idea FOKKO-0001-SUBS); esta idea cubre lo que falta.

Esta idea cubre solo: Log de peticiones sin datos personales. El resto del trabajo original está en las ideas hermanas `FOKKO-0062*-SUBS`.

## Alcance
- Log de peticiones sin datos personales

## Fuera
- Sesión, cookies y CSRF → FOKKO-0064-SUBS y FOKKO-0065A-SUBS a FOKKO-0065C-SUBS.

## Features afectadas
- `internal/shared` (`httpserver`)
- `cmd/api`

## Dudas abiertas
- Antes de escribir las tareas se revisa qué de esto ya existe en `router.go` y `run.go`, para no duplicarlo.

## Tareas
- [x] [TASK-001 Log de peticiones sin datos personales](TASK-001-REQUEST-LOGGING.md)
