# FOKKO-0033-SUBS · Override de monto y día

## Problema
El usuario puede reemplazar el monto o el día de cobro detectados de una línea. Cambia lo que se muestra, el total y el próximo cobro, pero no los cobros reales (decisión 10).

## Depende de
- FOKKO-0021-SUBS y FOKKO-0023-SUBS.

## Alcance
- Dominio `LineOverride` con monto y día opcionales, puerto, caso de uso y su efecto en el resumen.

## Fuera
- Endpoint → FOKKO-0074A-SUBS y FOKKO-0074B-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`, `ports`, `usecases`)

## Dudas abiertas
- Cómo se quita un override: se propone que guardar uno con monto y día vacíos lo borre, igual que el alias.

## Tareas
- [x] [TASK-001A Aceptar un override válido](TASK-001A-OVERRIDE-VALID.md)
- [x] [TASK-001B Rechazar un override inválido](TASK-001B-OVERRIDE-INVALID.md)
- [x] [TASK-001C Reconocer un override vacío](TASK-001C-OVERRIDE-EMPTY.md)
- [x] [TASK-002 Guardar o quitar un override](TASK-002-SAVE-OVERRIDE.md)
- [x] [TASK-003 El override cambia lo que se muestra](TASK-003-OVERRIDE-IN-SUMMARY.md)
