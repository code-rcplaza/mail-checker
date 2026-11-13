# FOKKO-0028-SUBS/TASK-002 · Descartar el aviso de un alza
Feature: `internal/subscriptions/usecases`

Como usuario quiero decir "ya lo vi, lo mantengo", para que deje de avisarme.

## Criterios
- Dada una línea con alza vigente, cuando se descarta, entonces se guarda con el monto nuevo.
- Dada una línea sin alza vigente, cuando se descarta, entonces devuelve `ErrNotFound`.
- Dada la línea de otro usuario, cuando se descarta, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
