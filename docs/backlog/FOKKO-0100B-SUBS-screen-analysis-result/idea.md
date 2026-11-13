# FOKKO-0100B-SUBS · Pantalla Análisis

## Problema
Mientras se sube, Cargando muestra el estado pendiente; después, Análisis resume qué se leyó y cuántas suscripciones son "seguras" y "por confirmar" (de `GET /api/summary`).

Esta idea cubre solo: Pantalla Análisis. El resto del trabajo original está en las ideas hermanas `FOKKO-0100*-SUBS`.

## Depende de
- FOKKO-0072A-SUBS a FOKKO-0072H-SUBS, FOKKO-0089A-SUBS y FOKKO-0089B-SUBS y FOKKO-0099-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Análisis
- Sin suscripciones encontradas

## Fuera
- Error → FOKKO-0101-SUBS.

## Capas
clean architecture por feature (`features/statements/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/statements`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Mostrar el archivo leído](TASK-001A-ANALYSIS-FILE.md)
- [ ] [TASK-001B Mostrar los dos conteos del análisis](TASK-001B-ANALYSIS-COUNTS.md)
- [ ] [TASK-001C Botón Continuar](TASK-001C-ANALYSIS-CONTINUE.md)
- [ ] [TASK-001D Texto de privacidad](TASK-001D-ANALYSIS-PRIVACY.md)
- [ ] [TASK-002 Sin suscripciones encontradas](TASK-002-ANALYSIS-EMPTY.md)
- [ ] [TASK-003 Conectar el hook con los componentes](TASK-003-CONTAINER-ANALYSIS.md)
