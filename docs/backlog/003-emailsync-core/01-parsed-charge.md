# 003-01 · `ParsedCharge` y su conversión a `Charge`
Feature: `internal/shared/domain`

Como desarrollador quiero convertir lo que extrae un parser en un `Charge` válido, con
una sola función, para que correo y cartola no dupliquen esa conversión.

## Criterios
- Dado un `ParsedCharge` con fecha y un origen (banco, fuente, external id, cuenta,
  fecha de respaldo), cuando se llama `ToCharge`, entonces devuelve un `Charge` con
  monto, moneda, comercio, tarjeta y fecha del `ParsedCharge`, y banco, fuente,
  external id y cuenta del origen.
- Dado un `ParsedCharge` con fecha cero, cuando se llama `ToCharge`, entonces el
  `Charge` usa la fecha de respaldo del origen.
- Dado un origen con external id vacío o solo espacios, cuando se llama `ToCharge`,
  entonces `ErrInvalid` con `common.external_id_missing`.
- Dado un `ParsedCharge` con monto negativo, cuando se llama `ToCharge`, entonces
  `ErrInvalid` con `common.amount_negative` (la validación es la de `NewCharge`).

## Fuera
- Origen desde un correo (Message-ID, fecha del correo) → 003-02.
- Origen desde una cartola → idea de `statements`.

## Notas
- Depende de 002-01 (`Charge`, `NewCharge`, `Money`).
- Datos sintéticos ("STREAMING CO").
