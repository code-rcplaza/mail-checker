# 004 · Núcleo de cartolas (fase 2, parte 3)

## Problema
Los avisos por correo no cubren todo: hay cobros que el banco nunca avisa y solo
aparecen en la cartola. Falta el flujo que importa una cartola, detecta de qué banco es
y la convierte en cobros, para que la vista "sin aviso" tenga datos. Esta idea lo
implementa con TDD y dobles en memoria, sin lectores de archivo reales ni base.

## Alcance
- `statements/domain`:
  - `Statement` (`ID`, `PublicID`, banco, nombre de archivo, `FileKind`, fecha de
    facturación, fecha de importación, cantidad de cobros, líneas malas) y `FileKind`
    (`pdf`, `xlsx`, `csv`, `text`).
  - Conversión de línea parseada a `Charge` con `ParsedCharge.ToCharge` (003-01):
    fuente `statement`, external id = `"<statementPublicID>:<n° de línea>"`, fecha de
    respaldo = fecha de facturación.
  - Llaves: `statements.bank_not_detected` (params `count`), `statements.bank_unknown`,
    `statements.already_imported`, `statements.no_charges` (params `badLines`),
    `statements.billing_date_missing`, `statements.file_empty`, `statements.not_found`,
    `statements.file_kind_unknown`.
- `statements/ports`: `FileExtractor`, `BankStatementParser`, `StatementRepository`
  (dobles en memoria).
- `statements/usecases`:
  - `ImportStatement(file, billingDate, bank?)`, en los pasos de `ARCHITECTURE.md`:
    1. extraer el texto;
    2. elegir banco: el indicado (si no está registrado → `bank_unknown`), o el único
       parser cuyo `Detect` lo reconoce (0 o más de 1 → `bank_not_detected`);
    3. parsear; las líneas sin fecha toman la fecha de facturación;
    4. guardar cartola y cobros en una transacción (`Create`).
    - Archivo vacío o sin fecha de facturación → `ErrInvalid`.
    - Ninguna línea válida → `ErrInvalid` con `no_charges`; no se guarda nada.
    - Ya existe una cartola del mismo banco con la misma fecha de facturación →
      `ErrConflict` con `already_imported`. Para reemplazarla se borra y se vuelve a
      subir.
    - `PublicID` lo genera un generador inyectado (como el reloj).
  - `ListStatements`, `RemoveStatement(publicID)` (`ErrNotFound`; borra también sus
    cobros), `ListBanks`.

## Fuera
- Lectores reales de PDF, Excel y CSV (`FileExtractor`) → fase 3.
- Parser real de la cartola BCI → fase 3 (necesita un PDF real anonimizado).
- Repositorio SQLite y generador UUID v7 real → fase 4.
- Endpoints `/api/statements` y `/api/banks` → fase 5.

## Features afectadas
- `internal/statements`

## Dudas abiertas
- Ninguna. Depende de 003-01 y 002-01. Decidido: una cartola por banco + fecha de facturación (reimportar →
  conflicto) y una cartola sin cobros se rechaza.
- Cambios aplicados a `ARCHITECTURE.md`:
  - external id de cartola: `"<statementPublicID>:<línea>"` en vez de
    `"<statementID>:<línea>"`, porque el `ID` interno lo asigna el repositorio dentro de
    `Create` y el cobro necesita su external id antes de guardarse.
  - tabla `statements`: `UNIQUE(bank, billing_date)`.
  - `ImportStatement`: agregar las reglas de conflicto y de cartola sin cobros.

## Historias
- [x] [01 Tipos de `statements`](01-statement-types.md)
- [x] [02 `ImportStatement`: camino feliz con banco elegido](02-import-statement-happy-path.md)
- [x] [03 `ImportStatement`: elegir banco](03-import-statement-bank-selection.md)
- [ ] [04 `ImportStatement`: validaciones](04-import-statement-guards.md)
- [ ] [05 `ListStatements`, `RemoveStatement`, `ListBanks`](05-list-remove-banks.md)
