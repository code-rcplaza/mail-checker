# FOKKO-0020B-SUBS/TASK-002 · Casos de uso de `statements` con usuario
Feature: `internal/statements/usecases`

Como usuario quiero que subir una cartola la deje a mi nombre, para verla solo yo.

## Criterios
- Dado `ImportStatement` con un usuario, cuando se importa, entonces la cartola y todos sus cargos quedan con ese `UserID`.
- Dado el mismo archivo subido por dos usuarios, cuando se importa, entonces no hay conflicto entre ellos (la unicidad es por usuario).
- Dado el mismo archivo subido dos veces por el mismo usuario, cuando se importa, entonces devuelve `ErrConflict` con `statements.already_imported`.

## Fuera
- Nada.

## Notas
- Ninguna.
