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
- Resuelta: los conteos salen de `GET /api/summary` por un método propio del gateway de statements
  (`readFound`, esquema mínimo): statements no importa subscriptions. Son de toda la cuenta.
- Resuelta: tras recargar, el archivo leído ya no está (no hay `GET` de una cartola); la pantalla
  muestra los conteos sin la parte del archivo.
- Resuelta: "Guardamos" sigue el esquema del back (incluye los últimos 4 dígitos de la tarjeta), no
  el texto del handoff; sale "Quitamos tus datos personales", que no era exacto.

## Tareas
- [x] [TASK-001A Mostrar el archivo leído](TASK-001A-ANALYSIS-FILE.md)
- [x] [TASK-001B Mostrar los dos conteos del análisis](TASK-001B-ANALYSIS-COUNTS.md)
- [x] [TASK-001C Botón Continuar](TASK-001C-ANALYSIS-CONTINUE.md)
- [x] [TASK-001D Texto de privacidad](TASK-001D-ANALYSIS-PRIVACY.md)
- [x] [TASK-002 Sin suscripciones encontradas](TASK-002-ANALYSIS-EMPTY.md)
- [x] [TASK-003 Conectar el hook con los componentes](TASK-003-CONTAINER-ANALYSIS.md)
