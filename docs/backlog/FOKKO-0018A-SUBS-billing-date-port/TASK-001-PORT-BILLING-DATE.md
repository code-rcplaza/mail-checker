# FOKKO-0018A-SUBS/TASK-001 · `BillingDate` en el puerto del parser
Feature: `internal/statements/ports`

Como usuario quiero que cada banco diga la fecha de su cartola, para no escribirla yo.

## Criterios
- Dado el puerto con `BillingDate`, cuando un doble de parser tiene fecha, entonces la devuelve con `true`.
- Dado un parser sin fecha, cuando se consulta, entonces devuelve `false`.

## Fuera
- Nada.

## Notas
- Hay que actualizar los dobles de prueba existentes de `BankStatementParser`.
