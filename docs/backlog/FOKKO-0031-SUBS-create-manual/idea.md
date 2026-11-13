# FOKKO-0031-SUBS · Crear suscripción manual

## Problema
El usuario puede agregar a mano una suscripción que no aparece en sus cartolas. Es una línea sin cobros que suma al total y cuyo próximo cobro es la siguiente vez que ocurre su día.

## Depende de
- FOKKO-0021-SUBS.

## Alcance
- Dominio `ManualSubscription` (nombre, monto, moneda, día de cobro) con validación.
- Puerto, caso de uso para crear y su aparición en el resumen.

## Fuera
- Editar y quitar → FOKKO-0032-SUBS.
- Endpoint → FOKKO-0075-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`, `ports`, `usecases`)

## Dudas abiertas
- Moneda permitida y largo máximo del nombre: el handoff las deja abiertas. Se propone CLP y USD, y 40 caracteres como el alias.

## Tareas
- [x] [TASK-001A Aceptar una manual válida](TASK-001A-MANUAL-VALID.md)
- [x] [TASK-001B Rechazar nombre o monto inválidos](TASK-001B-MANUAL-NAME-AMOUNT.md)
- [x] [TASK-001C Rechazar día o moneda inválidos](TASK-001C-MANUAL-DAY-CURRENCY.md)
- [x] [TASK-002 Próximo cobro de una manual](TASK-002-MANUAL-NEXT-CHARGE.md)
- [x] [TASK-003 Crear una suscripción manual](TASK-003-CREATE-MANUAL.md)
- [x] [TASK-004 La manual aparece en el resumen](TASK-004-MANUAL-IN-SUMMARY.md)
