# FOKKO-0100A-SUBS · Pantalla Cargando

## Problema
Mientras se sube, Cargando muestra el estado pendiente; después, Análisis resume qué se leyó y cuántas suscripciones son "seguras" y "por confirmar" (de `GET /api/summary`).

Esta idea cubre solo: Pantalla Cargando. El resto del trabajo original está en las ideas hermanas `FOKKO-0100*-SUBS`.

## Depende de
- FOKKO-0072A-SUBS a FOKKO-0072H-SUBS, FOKKO-0089A-SUBS y FOKKO-0089B-SUBS y FOKKO-0099-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Cargando

## Fuera
- Error → FOKKO-0101-SUBS.

## Capas
clean architecture por feature (`features/statements/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/statements`)

## Dudas abiertas
- Resuelta: Cargando es el estado pendiente de la misma pantalla `/fuente` (sin ruta propia); al
  responder navega a `/analisis/<id>`, y si falla vuelve al formulario con el mensaje.
- Resuelta: sin "Cancelar y elegir otra fuente" ni "Puedes salir y te avisamos". La subida no se
  interrumpe en el servidor (la cartola quedaría importada igual) y el MVP no tiene avisos.

## Tareas
- [x] [TASK-001 Cargando](TASK-001-UPLOAD-LOADING.md)
- [x] [TASK-002 Conectar el hook con los componentes](TASK-002-CONTAINER-LOADING.md)
