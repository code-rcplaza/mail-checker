---
name: story
description: Implementa una historia del backlog con TDD (plan → RED → GREEN → REFACTOR → VERIFY), con paradas para revisión. Tercer paso del loop.
argument-hint: "<NNN-MM>"
disable-model-invocation: true
---

# /story

Historia: $ARGUMENTS (formato `NNN-MM`)

Leer primero: la historia en `docs/backlog/NNN-*/MM-*.md`, su `idea.md`,
`docs/ARCHITECTURE.md` y `CLAUDE.md`. Si la historia no existe o ya está marcada
`[x]` en `idea.md`, parar y avisar.

## 1. Plan ⏸
- Listar los archivos a crear o modificar y qué se hace en cada uno.
- Listar los tests: un test (o caso de tabla) por criterio.
- Si falta un ejemplo real, pedirlo aquí.
- **Parar** y esperar el ok.

## 2. RED ⏸
- Escribir solo los tests, con los nombres del criterio que cubren.
- Correrlos y mostrar la salida: deben fallar **por la razón correcta** (no por
  error de compilación ajeno a la historia; para que compile se permiten firmas
  vacías que devuelven el valor cero).
- **Parar** y esperar el ok: los tests son la especificación.

## 3. GREEN
- Escribir el código mínimo para que los tests pasen. Nada fuera del plan.

## 4. REFACTOR
- Aplicar `CLAUDE.md`: funciones pequeñas compuestas como pasos, nombres que dicen
  lo que son o hacen, inmutabilidad, sin duplicación (regla de tres).
- Los tests siguen en verde después de cada cambio.

## 5. VERIFY ⏸
- Correr los comandos de verificación de `CLAUDE.md` (build, vet, staticcheck,
  tests, lint, tsc según corresponda) y **mostrar la salida**.
- Marcar la historia `[x]` en `idea.md`. Si era la última, marcar la idea `done`
  en `docs/backlog/README.md`.
- Si la historia cambió el diseño, proponer el cambio a `docs/ARCHITECTURE.md`
  (no editarlo sin ok).
- Resumen: archivos cambiados, por qué, qué quedó fuera.
- **Parar.** No commitear salvo que se pida (mensaje en inglés, con `NNN-MM`).

## Reglas
- Nunca modificar, saltar ni borrar un test para que pase.
- Si algo falla dos veces por lo mismo, parar y explicar.
- Alcance cerrado: solo los archivos del plan; si hace falta otro, preguntar.
- Si el diff pasa de ~300 líneas, parar y proponer partir la historia.
