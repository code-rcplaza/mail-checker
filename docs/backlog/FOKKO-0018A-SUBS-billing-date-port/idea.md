# FOKKO-0018A-SUBS · Puerto `BillingDate` en `BankStatementParser`

## Problema
Hoy `ImportStatement` exige la fecha de facturación de entrada. Debe salir del archivo (decisión 16) y pedirse al usuario solo si el parser no la encuentra.

Esta idea cubre solo: Puerto `BillingDate` en `BankStatementParser`. El resto del trabajo original está en las ideas hermanas `FOKKO-0018*-SUBS`.

## Alcance
- `BillingDate` en el puerto del parser

## Fuera
- La lectura concreta en BCI → FOKKO-0019-SUBS.

## Features afectadas
- `internal/statements` (`ports`, `usecases`)

## Dudas abiertas
- Si el usuario manda una fecha y el archivo trae otra, gana la del archivo. ¿Se acepta?

## Tareas
- [x] [TASK-001 `BillingDate` en el puerto del parser](TASK-001-PORT-BILLING-DATE.md)
