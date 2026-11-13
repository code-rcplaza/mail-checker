# FOKKO-0051B-SUBS/TASK-002B · Contar los envíos de la última hora
Feature: `internal/auth/adapters/sqlite`



## Criterios
- Dados tres envíos al mismo correo en la última hora, cuando se cuentan desde hace una hora, entonces son 3.
- Dado un envío de hace dos horas, cuando se cuentan desde hace una hora, entonces no cuenta.
- Dados envíos a correos distintos, cuando se cuenta uno, entonces solo cuentan los suyos.

## Fuera
- Nada.

## Notas
- Ninguna.
