# FOKKO-0002-SUBS/TASK-002 · Agrupar cobros
Feature: `internal/subscriptions/domain`

Como usuario quiero que los cobros del mismo servicio queden juntos aunque el banco
los describa distinto, para ver una suscripción por servicio.

## Criterios
- Dados cobros con descriptores distintos que el catálogo identifica como el mismo
  servicio conocido (y en monedas distintas), cuando se agrupan, entonces quedan en un
  solo grupo.
- Dados cobros de un servicio *umbrella* de 2.990 y 9.990 CLP, cuando se agrupan,
  entonces quedan en dos grupos (uno por monto exacto).
- Dados cobros *umbrella* del mismo monto en CLP y en USD, cuando se agrupan,
  entonces quedan en dos grupos.
- Dados cobros desconocidos cuyo nombre normalizado es igual, cuando se agrupan,
  entonces quedan juntos; con el mismo nombre pero otra moneda, quedan separados.
- Dado cualquier grupo, cuando se revisan sus cobros, entonces están ordenados por
  fecha ascendente, y el slice de entrada no fue modificado.

## Fuera
- Catálogo real → fase 3 (aquí `identify` y `normalize` son funciones falsas).

## Notas
- `domain` no importa `ports`: recibe `identify func(string) (Service, bool)` y
  `normalize func(string) string`; el caso de uso se las pasa desde `ServiceCatalog`.
