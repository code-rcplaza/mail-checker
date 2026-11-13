# FOKKO-0040-SUBS/TASK-002A · Crear un usuario
Feature: `internal/auth/usecases`



## Criterios
- Dados datos válidos, cuando se crea, entonces se guarda con `PublicID` UUID v7 generado por el caso de uso y la fecha de aceptación de términos con el reloj inyectado.
- Dados datos inválidos, cuando se crea, entonces no se llama al repositorio.

## Fuera
- Nada.

## Notas
- Ninguna.
