# 002-01 · Tipos de cobro
Feature: `internal/shared/domain`

Como desarrollador quiero tipos de cobro que no admitan valores inválidos para que
ninguna feature tenga que revalidar montos, monedas ni fechas.

## Criterios
- Dado `ParseCurrency("CLP")` o `ParseCurrency("USD")`, cuando se llama, entonces
  devuelve la moneda sin error.
- Dado `ParseCurrency("EUR")`, cuando se llama, entonces error `ErrInvalid` con llave
  `common.currency_unknown`.
- Dados datos válidos (banco, fuente, external id, fecha, monto, moneda, comercio,
  tarjeta), cuando se llama `NewCharge`, entonces devuelve un `Charge` con esos valores.
- Dado un monto negativo, cuando se llama `NewCharge`, entonces `ErrInvalid` con
  `common.amount_negative`.
- Dada una fecha cero, cuando se llama `NewCharge`, entonces `ErrInvalid` con
  `common.date_missing`.
- Dado un comercio vacío o solo espacios, cuando se llama `NewCharge`, entonces
  `ErrInvalid` con `common.merchant_missing`.

## Fuera
- Parseo de montos desde texto ("$7.500", "USD 20,00") → fase 3.
- `ParsedCharge` → idea de `emailsync`.

## Notas
- `Money` = centésimas en `int64` (CLP 5.990 → 599000).
- `Charge` se pasa por valor; `AccountID`/`StatementID` opcionales como punteros.
