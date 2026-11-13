# FOKKO-0050-SUBS/TASK-002 · Crear usuarios
Feature: `internal/auth/adapters/sqlite`

Como usuario quiero que mi registro se guarde una sola vez, para no tener cuentas duplicadas.

## Criterios
- Dado un usuario válido, cuando se crea, entonces se asigna el `ID` y se conserva el `PublicID` generado en Go.
- Dado un correo ya existente, cuando se crea, entonces devuelve `ErrConflict`.
- Dadas las fechas del usuario, cuando se guardan y se leen, entonces vuelven iguales (unix en segundos, UTC).

## Fuera
- Nada.

## Notas
- Ninguna.
