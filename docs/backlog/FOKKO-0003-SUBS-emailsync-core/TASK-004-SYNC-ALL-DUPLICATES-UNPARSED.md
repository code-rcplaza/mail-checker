# FOKKO-0003-SUBS/TASK-004 · `SyncAll`: duplicados y no reconocidos
Feature: `internal/emailsync/usecases`

Como usuario quiero que un aviso no se cuente dos veces y ver los correos que la app no
entendió.

## Criterios
- Dado el mismo aviso (mismo Message-ID) en dos cuentas, cuando se sincroniza, entonces
  se inserta una vez: `New` = 1 en la primera cuenta y 0 en la segunda; `Found` = 1 en
  ambas.
- Dado un aviso ya guardado en una sincronización anterior, cuando se vuelve a
  sincronizar, entonces `New` = 0.
- Dado un correo que el parser no reconoce, cuando se sincroniza, entonces se guarda
  con `UnparsedRepository.Save` y suma 1 a `Unparsed`; no suma a `New`.
- Dado un correo guardado como no reconocido que ahora se reconoce, cuando se
  sincroniza, entonces se inserta como cobro y se llama `Forget` con su Message-ID.
- Dado un aviso reconocido cuyo cobro es inválido (ej. comercio vacío), cuando se
  sincroniza, entonces se guarda como no reconocido y suma a `Unparsed`.

## Fuera
- Errores de lectura o de base → FOKKO-0003-SUBS/TASK-005.

## Notas
- La idempotencia la da `ChargeWriter.Insert` (devuelve `inserted`).
