# FOKKO-0072H-SUBS/TASK-001 · Próximos cobros
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver qué cobros vienen y cuándo, para llegar con saldo.

## Criterios
- Dados varios próximos cobros, cuando se serializa, entonces vienen ordenados por fecha con nombre, monto y `lineKey`.
- Dado un próximo cobro de una manual, cuando se serializa, entonces aparece igual que los demás.
- Dado ninguno, cuando se serializa, entonces la lista viene vacía y no `null`.

## Fuera
- Nada.

## Notas
- Ninguna.
