# FOKKO-0035-SUBS · Marcar "por cancelar"

## Problema
Sin lectura de correo no hay forma de saber si el usuario canceló. Él marca "por cancelar" a mano; la app no cancela nada (decisión 3).

## Depende de
- FOKKO-0021-SUBS.

## Alcance
- Puerto `Cancellations`, caso de uso para marcar y la marca en el resumen.

## Fuera
- Endpoint → FOKKO-0077A-SUBS y FOKKO-0077B-SUBS.
- Desmarcar: post-MVP.

## Features afectadas
- `internal/subscriptions` (`ports`, `usecases`, `domain`)

## Dudas abiertas
- ¿Una línea "por cancelar" sigue sumando al total hasta que deje de cobrarse? Se propone que sí, porque la app no sabe si se canceló.

## Tareas
- [ ] [TASK-001 Marcar "por cancelar"](TASK-001-MARK-FOR-CANCELLATION.md)
- [ ] [TASK-002 El resumen muestra la marca](TASK-002-SUMMARY-FLAGS-CANCELLATION.md)
