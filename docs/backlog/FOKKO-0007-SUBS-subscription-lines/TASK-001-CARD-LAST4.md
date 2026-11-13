# FOKKO-0007-SUBS/TASK-001 · Últimos 4 dígitos de la tarjeta
Feature: `internal/shared/domain`

Como usuario quiero que cada cobro recuerde con qué tarjeta se pagó para separar
suscripciones del mismo servicio pagadas con tarjetas distintas.

## Criterios
- Dado un `ParsedCharge` con `CardLast4 = "4821"`, cuando se llama `ToCharge`,
  entonces el `Charge` tiene `CardLast4 = "4821"`.
- Dado un `ParsedCharge` sin últimos 4 dígitos, cuando se llama `ToCharge`, entonces el
  `Charge` tiene `CardLast4 = ""` y no hay error (el aviso o la cartola puede no
  traerlos).
- Dados `"482"`, `"48211"`, `"48a1"` y `" 4821"`, cuando se construye el cobro
  (`NewCharge` / `ToCharge`), entonces `ErrInvalid` con
  `common.card_last4_invalid`.

## Fuera
- Extraer los dígitos de avisos y cartolas → post-MVP (avisos) y FOKKO-0011A-SUBS y FOKKO-0011B-SUBS (cartola).
- Columna `card_last4` en la base → fase 4.

## Notas
- Ejemplos de dígitos inventados.
