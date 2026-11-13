# FOKKO-0003-SUBS/TASK-006 · `SyncAll`: validaciones
Feature: `internal/emailsync/usecases`

Como usuario quiero que no se lancen dos sincronizaciones a la vez ni con un rango
inválido.

## Criterios
- Dado `months` = 0 o 25, cuando se llama `SyncAll`, entonces `ErrInvalid` con
  `emailsync.months_out_of_range` y no se lee ninguna cuenta.
- Dado `months` = 1 o 24, cuando se llama `SyncAll`, entonces sincroniza.
- Dada una sincronización en curso, cuando se llama otra `SyncAll`, entonces
  `ErrConflict` con `emailsync.sync_in_progress` sin leer cuentas.
- Dada una sincronización que terminó (con o sin error), cuando se llama otra, entonces
  sincroniza.

## Fuera
- Endpoint `/api/sync` → fase 5.

## Notas
- Rango 1–24 como constantes con nombre, igual que el resumen.
- El test de concurrencia usa un lector doble que bloquea hasta que el test lo libera.
