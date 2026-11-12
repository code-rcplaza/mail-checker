# 006-03 · `LooksLikeSubscription`
Feature: `internal/subscriptions/adapters/catalog`

Como usuario quiero que un único cobro con nombre de suscripción aparezca como posible,
sin esperar al segundo mes.

## Criterios
- Dados descriptores con `SUBSCRIPTION`, `MEMBERSHIP`, `SUSCRIPCION`, `MEMBRESIA`,
  `PREMIUM` o `PLAN` como palabra completa (ej. `"MICROSOFT. SUBSCRIPTION"`,
  `"Suscripción Anual"`, `"ACME Premium"`), cuando se llama, entonces `true`.
- Dados `"PLANETA"`, `"SUPREMIUM"`, `"MEMBERSHIPS"` y `"NETFLIX.COM"`, cuando se llama,
  entonces `false` (solo palabra completa).
- Dado `"MICROSOFT.SUBSCRIPTION"` (pegado con punto), cuando se llama, entonces `true`
  (la puntuación separa palabras al normalizar).

## Fuera
- Cómo se usa en la detección → ya cubierto en 002-07.

## Notas
- Se evalúa sobre el descriptor normalizado (006-01), así los acentos y mayúsculas no
  importan.
