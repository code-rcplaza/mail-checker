---
name: idea
description: Convierte una idea (o una fase de docs/ARCHITECTURE.md) en docs/backlog/NNN-slug/idea.md. Primer paso del loop idea → stories → story.
argument-hint: "<texto de la idea o número de fase>"
disable-model-invocation: true
---

# /idea

Entrada: $ARGUMENTS

## Pasos
1. Leer `docs/ARCHITECTURE.md` y `docs/backlog/README.md`.
2. Si la entrada es una fase de `ARCHITECTURE.md`, partir de su descripción.
3. Hacer las preguntas necesarias para cerrar alcance (no asumir). Esperar respuestas.
4. Elegir el siguiente número libre `NNN` (3 dígitos) y un slug corto en inglés,
   kebab-case. Proponerlo junto con el contenido.
5. Crear `docs/backlog/NNN-slug/idea.md` con la plantilla de abajo.
6. Agregar la fila en `docs/backlog/README.md` con estado `todo`.
7. **Parar.** Mostrar la idea y esperar aprobación. No crear historias.

## Plantilla `idea.md`
```markdown
# NNN · Título

## Problema
Qué falta y por qué importa (2–4 líneas).

## Alcance
- Lo que entra.

## Fuera
- Lo que no entra (y dónde iría, si se sabe).

## Features afectadas
- `internal/<feature>` …

## Dudas abiertas
- …

## Historias
(se completa con /stories)
```

## Reglas
- Contenido en español; slug en inglés.
- Nada de código ni de tests en este paso.
- Si la idea contradice `ARCHITECTURE.md`, decirlo y preguntar antes de seguir.
