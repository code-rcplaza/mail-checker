# FOKKO-0001-SUBS/TASK-001 · Módulo, errores de dominio y reloj
Feature: `backend/go.mod`, `internal/shared/domain`, `internal/shared/clock`

Como desarrollador quiero el módulo Go con los errores de dominio y un reloj
inyectable para que todas las features hablen el mismo idioma de errores y de tiempo.

## Criterios
- Dado `Invalid("statements.bank_not_detected", params)`, cuando se evalúa
  `errors.Is(err, ErrInvalid)`, entonces es `true`, y al extraer el error de dominio
  (`errors.As`) se obtienen la llave y los params.
- Dado `NotFound(key)` y `Conflict(key)`, cuando se evalúa `errors.Is`, entonces
  envuelven `ErrNotFound` y `ErrConflict` y conservan su llave.
- Dado un error de dominio envuelto con `fmt.Errorf("...: %w", err)`, cuando se
  evalúa `errors.Is` / `errors.As`, entonces siguen funcionando y la llave se conserva.
- Dadas las llaves comunes (`common.invalid_json`, `common.not_found`,
  `common.conflict`, `common.route_not_found`, `common.validation_failed`,
  `common.internal`), cuando se revisan, entonces existen como constantes y siguen
  el formato `<feature>.<motivo>` en snake_case.
- Dado `ErrNotFound`, cuando se compara con `ErrConflict` o `ErrInvalid`, entonces
  `errors.Is` es `false` (son distintos).
- Dado `clock.Now()`, cuando se llama, entonces devuelve la hora actual en UTC
  (entre un `time.Now()` previo y uno posterior, `Location() == time.UTC`).

## Fuera
- `Money`, `Charge` y demás tipos (fase 2).
- Errores por campo (`fields`): llegan con el primer formulario que los necesite.

## Notas
- El reloj es una función `func() time.Time` (regla de `CLAUDE.md`), no una interfaz.
