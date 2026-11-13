# FOKKO-0019-SUBS/TASK-001 · Leer la fecha del estado de cuenta
Feature: `internal/statements/adapters/bci`

Como usuario quiero que la app sepa de qué mes es mi cartola, para no equivocarme al escribirla.

## Criterios
- Dada una cartola con `FECHA ESTADO DE CUENTA`, cuando se llama `BillingDate`, entonces devuelve esa fecha a las 00:00 de `America/Santiago`.
- Dadas dos cartolas con distinto día de cierre, cuando se llama `BillingDate`, entonces cada una devuelve su propia fecha.

## Fuera
- Nada.

## Notas
- Ninguna.
