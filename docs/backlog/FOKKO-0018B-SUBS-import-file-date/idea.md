# FOKKO-0018B-SUBS · `ImportStatement` toma la fecha del archivo

## Problema
Hoy `ImportStatement` exige la fecha de facturación de entrada. Debe salir del archivo (decisión 16) y pedirse al usuario solo si el parser no la encuentra.

Esta idea cubre solo: `ImportStatement` toma la fecha del archivo. El resto del trabajo original está en las ideas hermanas `FOKKO-0018*-SUBS`.

## Depende de
- FOKKO-0018A-SUBS.

## Alcance
- `ImportStatement` toma la fecha del archivo

## Fuera
- La lectura concreta en BCI → FOKKO-0019-SUBS.

## Features afectadas
- `internal/statements` (`ports`, `usecases`)

## Dudas abiertas
- Si el usuario manda una fecha y el archivo trae otra, gana la del archivo. ¿Se acepta?

## Tareas
- [x] [TASK-001A Usar la fecha del archivo](TASK-001A-FILE-DATE-WINS.md)
- [x] [TASK-001B Usar la fecha del usuario o rechazar](TASK-001B-USER-DATE-FALLBACK.md)
