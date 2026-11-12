---
name: stories
description: Parte una idea aprobada (docs/backlog/NNN-slug/idea.md) en historias chicas con criterios Dado/Cuando/Entonces. Segundo paso del loop.
argument-hint: "<NNN>"
disable-model-invocation: true
---

# /stories

Idea: $ARGUMENTS

## Pasos
1. Leer `docs/backlog/$ARGUMENTS-*/idea.md` y `docs/ARCHITECTURE.md`.
   Si la idea no existe o tiene dudas abiertas, parar y preguntar.
2. Proponer la lista de historias (título + una línea cada una), en orden de
   implementación. Cada historia:
   - toca **una sola feature** (o `shared`);
   - se puede hacer con un diff de menos de ~300 líneas;
   - entrega algo verificable por tests.
3. **Parar** y esperar aprobación de la lista (se puede reordenar, partir o quitar).
4. Con la lista aprobada, crear un archivo por historia:
   `docs/backlog/NNN-slug/MM-slug.md` (MM = 2 dígitos, slug en inglés).
5. Completar la sección `## Historias` de `idea.md` como checklist:
   `- [ ] [MM Título](MM-slug.md)`.
6. Cambiar el estado de la idea a `doing` en `docs/backlog/README.md`.
7. **Parar.** Mostrar el resumen y esperar aprobación.

## Plantilla de historia
```markdown
# NNN-MM · Título
Feature: `internal/<feature>/<capa>`

Como usuario quiero … para …

## Criterios
- Dado …, cuando …, entonces ….
- Dado …, cuando …, entonces ….

## Fuera
- … (historia NNN-XX, si aplica)

## Notas
- Ejemplos reales que hacen falta (se piden, no se inventan).
```

## Reglas
- Las historias **no llevan estado**: el estado vive solo en el checklist de `idea.md`.
- Cada criterio debe poder convertirse en un test. Si no se puede, reescribirlo.
- Si hace falta un ejemplo real (formato de banco, cartola), anotarlo en Notas y
  pedirlo; nunca inventarlo.
- Nada de código ni de tests en este paso.
