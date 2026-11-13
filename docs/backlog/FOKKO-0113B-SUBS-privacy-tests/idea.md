# FOKKO-0113B-SUBS · Pruebas de privacidad

## Problema
El registro exige aceptar términos y política de privacidad, y la pantalla Análisis promete "solo guardamos servicio, monto y fecha". Esa promesa solo vale si el proceso realmente descarta el resto. Decisión abierta 12 de `docs/mvp.md`.

Esta idea cubre solo: Pruebas de privacidad. El resto del trabajo original está en las ideas hermanas `FOKKO-0113*-SUBS`.

## Depende de
- FOKKO-0094-SUBS y FOKKO-0070-SUBS.

## Alcance
- No se guarda el texto de la cartola
- Los logs no filtran datos

## Fuera
- Asesoría legal: los textos los define y revisa el usuario.

## Features afectadas
- `frontend` (`features/auth`, páginas legales)
- `internal/statements` (pruebas)

## Dudas abiertas
- ⛔ Decisión 12: los textos legales no se inventan. Se necesitan el texto de los términos y de la política, y qué dato se conserva y por cuánto tiempo (nombre y correo del usuario; servicio, monto y fecha de cada cobro).

## Tareas
- [ ] [TASK-001 No se guarda el texto de la cartola](TASK-001-PRIVACY-NO-RAW-TEXT.md)
- [ ] [TASK-002 Los logs no filtran datos](TASK-002-PRIVACY-LOGS.md)
