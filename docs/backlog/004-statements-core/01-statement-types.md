# 004-01 · Tipos de `statements`
Feature: `internal/statements/domain`

Como desarrollador quiero los tipos de la cartola y la conversión de sus líneas a
cobros para que el caso de uso solo orqueste pasos.

## Criterios
- Dado `ParseFileKind("pdf" | "xlsx" | "csv" | "text")`, cuando se llama, entonces
  devuelve el tipo sin error.
- Dado `ParseFileKind("doc")`, cuando se llama, entonces `ErrInvalid` con
  `statements.file_kind_unknown`.
- Dadas líneas parseadas con fecha, cuando se llama
  `chargesFromLines(lines, bank, publicID, billingDate)`, entonces cada cobro tiene
  fuente `statement`, el banco, external id `"<publicID>:<n>"` (n desde 1, en el orden de
  las líneas) y su propia fecha.
- Dada una línea sin fecha, cuando se convierte, entonces el cobro toma la fecha de
  facturación.
- Dada una línea que no se puede convertir (ej. comercio vacío), cuando se convierte,
  entonces no genera cobro y se devuelve su número `n` en las líneas malas (el front
  arma el texto); las demás líneas se convierten igual.
- Las llaves `statements.bank_not_detected`, `statements.bank_unknown`,
  `statements.already_imported`, `statements.no_charges`,
  `statements.billing_date_missing`, `statements.file_empty`,
  `statements.not_found` y `statements.file_kind_unknown` existen como constantes.

## Fuera
- Uso en `ImportStatement` → 004-02.

## Notas
- Depende de 003-01 (`ParsedCharge.ToCharge`) y 002-01 (`Charge`).
- Agrega la llave `statements.file_kind_unknown`, que `ParseFileKind` necesita.
- Datos sintéticos ("STREAMING CO").
