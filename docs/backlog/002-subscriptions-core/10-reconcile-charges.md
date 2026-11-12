# 002-10 · `ReconcileCharges`
Feature: `internal/subscriptions/usecases`

Como usuario quiero ver qué cobros de la cartola nunca me avisaron por correo.

## Criterios
- Dado un cobro en cartola y en correo que coinciden, cuando se reconcilia, entonces
  queda en avisados.
- Dado un cobro solo en la cartola, cuando se reconcilia, entonces queda en sin aviso.
- Dado un cobro solo por correo, cuando se reconcilia, entonces queda en solo correo.
- Dados dos cobros iguales en cartola y uno en correo, cuando se reconcilia, entonces
  uno queda en avisados y el otro en sin aviso (1 a 1).
- Dado `months` fuera de 1–24 o un `ChargeReader` que falla, cuando se reconcilia,
  entonces el mismo error que en `GetSummary`.
