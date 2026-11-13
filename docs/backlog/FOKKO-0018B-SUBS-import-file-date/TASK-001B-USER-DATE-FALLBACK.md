# FOKKO-0018B-SUBS/TASK-001B · Usar la fecha del usuario o rechazar
Feature: `internal/statements/usecases`



## Criterios
- Dado un archivo sin fecha y con fecha del usuario, cuando se importa, entonces usa la del usuario.
- Dado un archivo sin fecha y sin fecha del usuario, cuando se importa, entonces devuelve `ErrInvalid` con `statements.billing_date_required` y no guarda nada.

## Fuera
- Nada.

## Notas
- Hoy `validateInput` rechaza la fecha en cero antes de leer el archivo; ese orden cambia.
