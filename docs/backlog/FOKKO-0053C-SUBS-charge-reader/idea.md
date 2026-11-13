# FOKKO-0053C-SUBS · Lector de cargos de `subscriptions`

## Problema
`subscriptions` necesita leer los cargos del usuario sin importar la feature `statements`.

## Depende de
- FOKKO-0021-SUBS y FOKKO-0053A-SUBS.

## Alcance
- Adaptador sqlite de `ChargeReader` en `subscriptions`.

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Leer los cargos del usuario](TASK-001-CHARGE-READER.md)
