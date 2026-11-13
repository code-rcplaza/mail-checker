---
name: stories
description: Parte una idea aprobada (docs/backlog/FOKKO-XXXX-SUBS-slug/idea.md) en tareas chicas con criterios Dado/Cuando/Entonces. Segundo paso del loop.
argument-hint: "<XXXX>"
disable-model-invocation: true
---

# /stories

Idea: $ARGUMENTS (4 dígitos, ej. `0006` = `FOKKO-0006-SUBS`)

## Pasos
1. Leer `docs/backlog/FOKKO-$ARGUMENTS-SUBS-*/idea.md` y `docs/ARCHITECTURE.md`.
   Si la idea no existe o tiene dudas abiertas, parar y preguntar.
2. Proponer la lista de tareas (título + una línea cada una), en orden de
   implementación. Cada tarea:
   - toca **una sola feature** (o `shared`);
   - se puede hacer con un diff de menos de ~300 líneas;
   - entrega algo verificable por tests.
3. **Parar** y esperar aprobación de la lista (se puede reordenar, partir o quitar).
4. Con la lista aprobada, crear un archivo por tarea:
   `docs/backlog/FOKKO-XXXX-SUBS-slug/TASK-XXX-NOMBRE.md` (XXX = 3 dígitos desde 001,
   reinicia en cada idea; NOMBRE en inglés, mayúsculas y guiones).
5. Completar la sección `## Tareas` de `idea.md` como checklist:
   `- [ ] [TASK-XXX Título](TASK-XXX-NOMBRE.md)`.
6. Cambiar el estado de la idea a `doing` en `docs/backlog/README.md`.
7. **Parar.** Mostrar el resumen y esperar aprobación.

## Plantilla de tarea
```markdown
# FOKKO-XXXX-SUBS/TASK-XXX · Título
Feature: `internal/<feature>/<capa>`

Como usuario quiero … para …

## Criterios
- Dado …, cuando …, entonces ….
- Dado …, cuando …, entonces ….

## Fuera
- … (tarea FOKKO-XXXX-SUBS/TASK-XXX, si aplica)

## Notas
- Ejemplos reales que hacen falta (se piden, no se inventan).
```

## Reglas
- Las tareas **no llevan estado**: el estado vive solo en el checklist de `idea.md`.
- Cada criterio debe poder convertirse en un test. Si no se puede, reescribirlo.
- Si hace falta un ejemplo real (formato de banco, cartola), anotarlo en Notas y
  pedirlo; nunca inventarlo.
- Nada de código ni de tests en este paso.
