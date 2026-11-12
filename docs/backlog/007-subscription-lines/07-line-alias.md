# 007-07 · Alias de línea
Feature: `internal/subscriptions/usecases` + `ports` + `domain`

Como usuario quiero ponerle un nombre a cada línea (Personal, Padres, Hijos) para
saber de quién es cada suscripción.

## Criterios
- Dada una línea existente (activa, posible o inactiva), cuando se llama
  `NameLine(key, "  Padres ")`, entonces se guarda `"Padres"` para esa clave.
- Dado un alias vacío o solo espacios, cuando se llama `NameLine`, entonces se borra el
  alias de esa línea.
- Dado un alias de más de 40 caracteres, cuando se llama, entonces `ErrInvalid` con
  `subscriptions.alias_too_long` y no se guarda.
- Dada una clave que no es una línea del período, cuando se llama, entonces
  `ErrNotFound` con `subscriptions.line_not_found`.
- Dados alias guardados, cuando se pide el resumen (`GetSummary`), entonces cada
  `Subscription` trae el alias de su clave (vacío si no tiene).
- Dado que `LineAliases` falla al leer o guardar, cuando se llama, entonces el error se
  devuelve envuelto.

## Fuera
- Endpoint y pantalla → fases 5 y 7.

## Notas
- Puerto nuevo `LineAliases {List, Save, Delete}` con doble en memoria.
- Si la clave de la línea cambia (cambia de día de verdad), el alias se pierde, igual
  que una decisión (aceptado con la clave b).
- El largo se cuenta en caracteres (runas), no en bytes.
