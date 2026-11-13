# FOKKO-0034-SUBS · "Dejar de seguir" una línea

## Problema
"Dejar de seguir" oculta una línea: no suma ni aparece. La marca se guarda por clave de línea para que no vuelva al recalcular (decisión 10).

## Depende de
- FOKKO-0021-SUBS.

## Alcance
- Caso de uso que oculta una línea y su efecto en el resumen.

## Fuera
- Endpoint → FOKKO-0078-SUBS.
- Recuperarla: se agrega de nuevo como manual.

## Features afectadas
- `internal/subscriptions` (`ports`, `usecases`, `domain`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Ocultar una línea](TASK-001-STOP-FOLLOWING.md)
- [ ] [TASK-002 El resumen no muestra la línea oculta](TASK-002-SUMMARY-HIDES-LINE.md)
