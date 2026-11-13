# FOKKO-0003-SUBS/TASK-005 · `SyncAll`: errores por cuenta
Feature: `internal/emailsync/usecases`

Como usuario quiero que si falla una cuenta las demás se sincronicen igual y saber cuál
falló.

## Criterios
- Dado un lector que devuelve `ErrMailAuthRejected` para una cuenta, cuando se
  sincroniza, entonces su `SyncResult` tiene `ErrorKey` = `emailsync.mail_auth_failed`
  y las demás cuentas se sincronizan.
- Dado un lector que devuelve otro error para una cuenta, cuando se sincroniza,
  entonces su `ErrorKey` = `emailsync.mail_read_failed` y las demás siguen.
- Dada una sincronización, cuando termina cada cuenta, entonces se llama `RecordSync`
  con la cuenta, la hora del reloj y su `ErrorKey` (vacía si salió bien).
- Dado un error de `Credentials`, `Insert`, `Save`, `Forget` o `RecordSync`, cuando se
  sincroniza, entonces `SyncAll` se corta y devuelve el error envuelto con contexto.

## Fuera
- Detección real de credenciales rechazadas en IMAP → fase 5.

## Notas
- Un error de base afecta a todas las cuentas: fail fast en vez de seguir.
