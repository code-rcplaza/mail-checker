# FOKKO-0006-SUBS/TASK-004 · Catálogo con la agrupación real
Feature: `internal/subscriptions/adapters/catalog`

Como usuario quiero que los cobros reales de mis cartolas queden en las líneas
correctas cuando el catálogo real se usa con la detección.

## Criterios
- Dado el catálogo, entonces cumple `ports.ServiceCatalog` (verificado al compilar).
- Dados tres cobros mensuales (mismo día y monto) con los descriptores
  `"NETFLIX.COM 866-579-7172"`, `"NETFLIX LAS CONDES"` y `"NETFLIX"`, cuando se agrupan
  con `GroupCharges` y el catálogo, entonces forman una sola línea cuyo padre es
  `service:Netflix`.
- Dados tres cobros mensuales de `"DL*GOOGLE YOUTUBE"` de CLP 21.700 y tres de
  CLP 11.000, cuando se agrupan, entonces forman dos líneas con el mismo padre
  `umbrella:Google`.
- Dados tres cobros mensuales de `"PAYU *UBER TRIP"` en CLP, cuando se agrupan, entonces
  forman una línea cuyo padre es el comercio desconocido `merchant:UBER TRIP:CLP`
  (el prefijo del procesador no cambia el padre).
- Dado un único cobro `"MICROSOFT. SUBSCRIPTION"`, cuando se agrupa, entonces cae en el
  grupo suelto `umbrella:Microsoft|one-off`, marcado como servicio conocido.

## Fuera
- Separar el lugar pegado a la descripción (`SANTIAGO`, teléfono) → parser de cartola
  (FOKKO-0011A-SUBS y FOKKO-0011B-SUBS, parser de cartola BCI).
- Claves de línea con tarjeta y día (`padre|tarjeta|dDD`) → ya cubiertas en FOKKO-0007-SUBS/TASK-004; aquí
  solo se verifica el padre y la cantidad de líneas.

## Notas
- Usa `domain.GroupCharges` real; nada de dobles.
- Fechas y montos exactos de cada caso se fijan al escribir el test; deben cumplir las
  reglas de cadencia de `ChainLines` (tres cobros mensuales bastan para formar línea).
- Reescrita tras la idea FOKKO-0007-SUBS: las claves ya no llevan monto ni moneda por grupo.
