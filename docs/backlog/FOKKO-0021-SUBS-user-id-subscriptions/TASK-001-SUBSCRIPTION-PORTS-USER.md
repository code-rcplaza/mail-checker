# FOKKO-0021-SUBS/TASK-001 · Puertos de `subscriptions` por usuario
Feature: `internal/subscriptions/ports`

Como usuario quiero que mis decisiones y alias sean solo míos, para que nadie más los vea.

## Criterios
- Dados `ChargeReader`, `SubscriptionDecisions` y `LineAliases`, cuando se llaman, entonces reciben el `userID`.
- Dados dos usuarios, cuando uno guarda una decisión con la misma clave de línea, entonces el otro no la ve.

## Fuera
- Nada.

## Notas
- Ninguna.
