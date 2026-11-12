# 002-11 · Confirmar o descartar una posible
Feature: `internal/subscriptions/usecases`

Como usuario quiero decir si una posible suscripción lo es o no, y si es mensual o
anual, para que cuente bien (o deje de aparecer) en mi gasto.

## Criterios
- Dada una posible actual, cuando se confirma como mensual o como anual, entonces la
  decisión `confirmed` queda guardada con su llave de grupo y esa cadencia.
- Dada una confirmación sin cadencia (o con una que no es `monthly` ni `yearly`),
  cuando se decide, entonces `ErrInvalid` con `subscriptions.cadence_required`.
- Dada una posible actual, cuando se descarta, entonces la decisión `dismissed` queda
  guardada.
- Dada una llave que no corresponde a una posible actual, cuando se decide, entonces
  `ErrNotFound` con `subscriptions.possible_not_found`.
- Dada una decisión que no es `confirm` ni `dismiss`, cuando se decide, entonces
  `ErrInvalid` con `subscriptions.decision_unknown`.
- Dado que `SubscriptionDecisions` falla al guardar, cuando se decide, entonces
  devuelve ese error envuelto.

## Notas
- La llave de grupo es estable entre recálculos (ej. `merchant:STREAMING CO:CLP`).
- Para saber si es "posible actual" se recalcula con la misma lógica de `GetSummary`
  (sin duplicarla: se reutiliza el paso de detección) sobre la ventana máxima (24 meses),
  para que una posible visible en cualquier período se pueda decidir.
