# FOKKO-0003-SUBS/TASK-003 · `SyncAll`: camino feliz
Feature: `internal/emailsync/usecases`

Como usuario quiero sincronizar mis cuentas para que los avisos de compra queden como
cobros.

## Criterios
- Dada una cuenta y un banco registrado, cuando se llama `SyncAll(ctx, 3)`, entonces el
  lector recibe el filtro del banco con "desde" = hoy − 3 meses (reloj inyectado).
- Dados avisos que el parser reconoce, cuando se sincroniza, entonces cada uno se
  inserta como `Charge` con fuente `email` y el `SyncResult` de la cuenta tiene
  `Found` = avisos leídos y `New` = insertados.
- Dadas dos cuentas, cuando se sincroniza, entonces hay un `SyncResult` por cuenta, en
  el orden del directorio, con su `AccountPublicID` y correo.
- Dados dos bancos registrados, cuando se sincroniza una cuenta, entonces se busca con
  el filtro de cada banco y cada aviso se parsea con el parser de su banco.

## Fuera
- Duplicados y no reconocidos → FOKKO-0003-SUBS/TASK-004.
- Errores → FOKKO-0003-SUBS/TASK-005. Validaciones y concurrencia → FOKKO-0003-SUBS/TASK-006.

## Notas
- Puertos `AccountDirectory`, `MailReader`, `BankEmailParser`, `ChargeWriter`,
  `UnparsedRepository` en `ports`; dobles en memoria en los tests.
- Avisos sintéticos; sin datos reales.
