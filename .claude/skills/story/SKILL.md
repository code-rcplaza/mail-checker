---
name: story
description: Implementa una historia del backlog con TDD (plan → RED → GREEN → REFACTOR → VERIFY) de corrido, con una sola parada al final con el informe. Tercer paso del loop.
argument-hint: "<NNN-MM>"
disable-model-invocation: true
---

# /story

Historia: $ARGUMENTS (formato `NNN-MM`)

Leer primero: la historia en `docs/backlog/NNN-*/MM-*.md`, su `idea.md`,
`docs/ARCHITECTURE.md` y `CLAUDE.md`. Si la historia no existe o ya está marcada
`[x]` en `idea.md`, parar y avisar.

El ciclo corre **completo sin parar** y termina con un informe. Solo se para antes
si aparece una de las **paradas obligatorias** (abajo).

## 1. Plan
- Definir los archivos a crear o modificar y qué se hace en cada uno.
- Definir los tests: un test (o caso de tabla) por criterio.

## 2. RED
- Escribir solo los tests, con los nombres del criterio que cubren.
- Para que compile se permiten stubs que devuelven el valor cero.
- Correrlos y **guardar la salida**: deben fallar por la razón correcta.
- Si un test pasa en RED (falso verde), reforzarlo hasta que falle por la razón
  correcta, sin cambiar lo que el criterio pide.

## 3. GREEN
- Escribir el código mínimo para que los tests pasen. Nada fuera del plan.

## 4. REFACTOR
- Aplicar `CLAUDE.md`: funciones pequeñas compuestas como pasos, nombres que dicen
  lo que son o hacen, inmutabilidad, sin duplicación (regla de tres).
- Los tests siguen en verde después de cada cambio.

## 5. VERIFY
- Correr los comandos de verificación de `CLAUDE.md` (build, vet, staticcheck,
  gofmt, tests, lint, tsc según corresponda).
- Si todo pasa: marcar la historia `[x]` en `idea.md`. Si era la última, marcar la
  idea `done` en `docs/backlog/README.md`.

## 6. Informe ⏸ (única parada normal)
- ✅ / ❌ al inicio: si la historia terminó correctamente.
- Plan aplicado: archivos y qué hace cada uno.
- RED: tabla test → por qué fallaba (de la salida guardada).
- Refactor: qué cambió (o "nada").
- VERIFY: salida de cada comando.
- Decisiones tomadas dentro del alcance y lo que quedó fuera.
- Si la historia cambió el diseño, proponer el cambio a `docs/ARCHITECTURE.md`
  (no editarlo sin ok).
- No commitear salvo que se pida (mensaje en inglés, con `NNN-MM`).

## Paradas obligatorias (parar, explicar y esperar)
- Hace falta una decisión del usuario: diseño, nombres que cambian
  `docs/ARCHITECTURE.md`, o una dependencia no aprobada.
- Falta un ejemplo real (formato de banco, cartola): se pide, no se inventa.
- El diff pasa de ~300 líneas: proponer partir la historia.
- Algo falla dos veces por lo mismo.
- Un test sigue en falso verde y no se puede reforzar sin cambiar el criterio.
- Hace falta tocar un archivo fuera del plan que no es de la historia.

## Reglas
- Nunca modificar, saltar ni borrar un test para que pase.
