# 006-04 · Catálogo con la agrupación real
Feature: `internal/subscriptions/adapters/catalog`

Como usuario quiero que los cobros reales de mis cartolas queden en los grupos
correctos cuando el catálogo real se usa con la detección.

## Criterios
- Dado el catálogo, entonces cumple `ports.ServiceCatalog` (verificado al compilar).
- Dados `"NETFLIX.COM 866-579-7172"` y `"NETFLIX LAS CONDES"` en CLP, cuando se agrupan
  con `GroupCharges` y el catálogo, entonces quedan en un solo grupo `service:Netflix`.
- Dados dos cobros `"DL*GOOGLE YOUTUBE"` de CLP 21.700 y CLP 11.000, cuando se agrupan,
  entonces quedan en dos grupos `umbrella:Google:2170000:CLP` y
  `umbrella:Google:1100000:CLP`.
- Dados dos cobros `"PAYU *UBER TRIP"` en CLP, cuando se agrupan, entonces quedan en
  un solo grupo de comercio desconocido `merchant:UBER TRIP:CLP`.
- Dado `"MICROSOFT. SUBSCRIPTION"` con un solo cobro, cuando se agrupa, entonces es un
  grupo umbrella de Microsoft.

## Fuera
- Separar el lugar pegado a la descripción (`SANTIAGO`, teléfono) → parser de cartola
  (idea 010).

## Notas
- Usa `domain.GroupCharges` real; nada de dobles.
