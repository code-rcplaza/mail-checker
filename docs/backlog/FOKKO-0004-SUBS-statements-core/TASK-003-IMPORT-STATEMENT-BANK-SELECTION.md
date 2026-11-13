# FOKKO-0004-SUBS/TASK-003 · `ImportStatement`: elegir banco
Feature: `internal/statements/usecases`

Como usuario quiero que la app detecte de qué banco es la cartola si no lo elijo.

## Criterios
- Dado un banco elegido que no está registrado, cuando se importa, entonces `ErrInvalid`
  con `statements.bank_unknown` y no se parsea.
- Dado que no se elige banco y exactamente un parser lo reconoce con `Detect`, cuando se
  importa, entonces se usa ese parser.
- Dado que no se elige banco y ningún parser lo reconoce, cuando se importa, entonces
  `ErrInvalid` con `statements.bank_not_detected` y `count` = 0.
- Dado que no se elige banco y dos parsers lo reconocen, cuando se importa, entonces
  `ErrInvalid` con `statements.bank_not_detected` y `count` = 2.
- Dado un banco elegido registrado, cuando se importa, entonces no se llama `Detect`.

## Fuera
- Validaciones de archivo y fecha → FOKKO-0004-SUBS/TASK-004.
