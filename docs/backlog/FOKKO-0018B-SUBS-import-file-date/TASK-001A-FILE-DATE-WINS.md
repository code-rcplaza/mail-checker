# FOKKO-0018B-SUBS/TASK-001A · Usar la fecha del archivo
Feature: `internal/statements/usecases`



## Criterios
- Dado un archivo con fecha y sin fecha del usuario, cuando se importa, entonces la cartola queda con la fecha del archivo.
- Dado un archivo con fecha y otra del usuario, cuando se importa, entonces gana la del archivo.

## Fuera
- Nada.

## Notas
- Hoy `validateInput` rechaza la fecha en cero antes de leer el archivo; ese orden cambia.
