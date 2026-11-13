# FOKKO-0061-SUBS · Limpiar códigos y sesiones vencidos

## Problema
Códigos, envíos y sesiones vencidos se acumulan sin límite.

## Depende de
- FOKKO-0051A-SUBS y FOKKO-0051B-SUBS y FOKKO-0052-SUBS.

## Alcance
- Caso de uso que borra lo vencido y consultas que lo hacen.

## Fuera
- Dónde y cada cuánto corre → FOKKO-0082A-SUBS a FOKKO-0082D-SUBS.

## Features afectadas
- `internal/auth` (`usecases`, `adapters/sqlite`)

## Dudas abiertas
- Se propone una tarea periódica que arranca `main.go` (cada hora, cancelable con el contexto). Lo vencido ya no vale aunque no se borre, así que no es urgente.

## Tareas
- [ ] [TASK-001 Borrar lo vencido](TASK-001-CLEANUP-QUERIES.md)
- [ ] [TASK-002 Caso de uso de limpieza](TASK-002-CLEANUP-USECASE.md)
