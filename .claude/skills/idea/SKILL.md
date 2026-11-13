---
name: idea
description: Crea el idea.md de una idea del plan (docs/backlog/README.md) como docs/backlog/FOKKO-XXXX-SUBS-slug/idea.md. Primer paso del loop idea → stories → story.
argument-hint: "<XXXX> [contexto opcional]"
disable-model-invocation: true
---

# /idea

Entrada: $ARGUMENTS (el primer token es el ID de 4 dígitos, ej. `0008` = `FOKKO-0008-SUBS`;
el resto es contexto opcional)

El plan manda: número y título salen de la fila de la idea en `docs/backlog/README.md`.
El texto extra solo aporta contexto para las preguntas; nunca cambia el título ni el número.

## Pasos
1. Leer `docs/backlog/README.md` y `docs/ARCHITECTURE.md`.
2. Buscar `FOKKO-XXXX-SUBS` en el plan. Parar y avisar si:
   - el ID no aparece en el plan (una idea nueva se agrega primero al plan, con el ok
     del usuario, y recién después se crea su `idea.md`);
   - la fila ya tiene enlace o existe `docs/backlog/FOKKO-XXXX-SUBS-*/`;
   - alguna idea de la que depende (número menor del mismo bloque que la bloquea) sigue
     en `todo`: decirlo y preguntar si se sigue igual.
3. Partir del título de la fila y del contexto de `ARCHITECTURE.md` y `docs/mvp.md`.
4. Hacer las preguntas necesarias para cerrar alcance (no asumir). Esperar respuestas.
5. Proponer un slug corto en inglés, kebab-case, junto con el contenido.
6. Crear `docs/backlog/FOKKO-XXXX-SUBS-slug/idea.md` con la plantilla de abajo.
7. **Actualizar la fila existente** en `docs/backlog/README.md`: el título pasa a ser
   un enlace al `idea.md` y el estado sigue en `todo`. No se agrega una fila nueva.
8. **Parar.** Mostrar la idea y esperar aprobación. No crear tareas.

## Plantilla `idea.md`
```markdown
# FOKKO-XXXX-SUBS · Título (el de la fila del plan)

## Problema
Qué falta y por qué importa (2–4 líneas).

## Alcance
- Lo que entra.

## Fuera
- Lo que no entra (y dónde iría, con su `FOKKO-XXXX-SUBS`, si se sabe).

## Features afectadas
- `internal/<feature>` …

## Dudas abiertas
- …

## Tareas
(se completa con /stories)
```

## Reglas
- Contenido en español; slug en inglés.
- Una idea es unitaria: un resultado verificable, una capa, una feature; si el título lleva "y", se parte con sufijo de letra (`FOKKO-0072A-SUBS`);
  si el título lleva una "y", se parte (y se avisa para ajustar el plan).
- Nada de código ni de tests en este paso.
- Si la idea contradice `ARCHITECTURE.md`, decirlo y preguntar antes de seguir.
- Si el alcance real no coincide con el título del plan, no se renombra en silencio:
  se propone el cambio al plan y se espera el ok.
