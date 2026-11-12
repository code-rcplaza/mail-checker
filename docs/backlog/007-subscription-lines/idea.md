# 007 · Líneas de suscripción (corrección de la detección)

## Problema
La detección supone que un servicio (o comercio) es **una sola** suscripción. Probado
con escenarios de cartolas reales, eso esconde suscripciones y gasto:
- dos Netflix (el propio y el de los padres) → 1 posible sin cadencia, gasto $0;
- dos suscripciones Apple de $5.990 (días 5 y 9) → 1 posible, gasto $0;
- dos pólizas `SEGUROS CLC` cobradas el mismo día → ocultas en "otros cobros".

Los cobros de varias suscripciones paralelas se mezclan, los intervalos se acortan y
la cadencia desaparece. Además, el umbrella separado por monto exacto parte una
suscripción cuando sube de precio. Como el usuario toma decisiones con esta
información, la detección debe separar cada suscripción aunque compartan servicio.

## Alcance
- **Modelo**: el grupo actual pasa a ser el **padre** (servicio, umbrella o comercio +
  moneda); dentro del padre, cada **línea** es una suscripción: la cadena de cobros
  mes a mes de una misma cosa. Cada cobro se guarda aparte; nada se sobreescribe.
  ```
  Netflix (padre)
  ├─ línea 1: cuenta A · tarjeta 4821 · día 01 · $12.990
  └─ línea 2: cuenta B · tarjeta 4821 · día 03 · $7.990
  ```
- `shared/domain`: `Charge` y `ParsedCharge` ganan los **últimos 4 dígitos de la
  tarjeta** (vacío si el aviso o la cartola no los trae; si viene, exactamente 4
  dígitos). `ToCharge` los copia.
- `subscriptions/domain`:
  1. **Origen**: dentro del padre, los cobros se separan primero por cuenta de correo y
     por últimos 4 dígitos cuando se conocen. Un cobro de cartola que se cruza con su
     aviso hereda la cuenta del aviso; uno sin aviso no tiene cuenta y se une a la
     línea con la que encadena.
  2. **Cadena**: dentro del mismo origen, un cobro continúa una línea si está a un
     intervalo válido de su último cobro (mensual k × 24–38 o anual 330–400 días).
     Con varias líneas posibles, va a la de monto más parecido y, si empatan, a la de
     fecha esperada más cercana. Si ninguna calza, abre una línea nueva.
  3. **Salto de monto** (> 25 %) en la fecha de una línea: en el mismo origen es
     **cambio de plan** dentro de la línea; en otro origen ya es otra línea por (1).
  4. **Umbrella** deja de separarse por monto exacto: sus líneas salen de (1)–(3).
  5. **Clave estable de la línea** (opción b): padre + últimos 4 dígitos + día de cobro
     típico (mediana de los días de la línea), ej. `service:Netflix|4821|d01`.
     Sin ids internos ni correos en la clave (sale por la API). Si dos líneas chocan
     en la misma clave, la segunda lleva `#2` (por orden de su primer cobro).
  6. Detección, monto actual, cambios de plan, gasto mensual, resumen y decisiones
     trabajan **por línea**; `Subscription` lleva su padre para que la vista agrupe.
  7. El cruce cartola ↔ correo exige además los mismos últimos 4 dígitos cuando
     ambos los traen.
  8. **Cobro suelto**: un cobro que no encadena, en el mismo origen que una línea
     activa, va a "otros cobros" del padre (como 002-12); en otro origen abre una
     línea nueva (posible).
  9. **Alias**: el usuario nombra cualquier línea (Personal, Padres…); se guarda por
     clave de línea en un puerto nuevo `LineAliases` y el resumen lo muestra.
- Tests con los escenarios reales de la auditoría: dos Netflix (mismo y distinto
  precio, misma y distinta cuenta), dos Apple de $5.990, dos pólizas el mismo día,
  Google con dos montos, cambio de plan de Netflix en la misma cuenta, y dos líneas
  iguales el mismo día (se separan, la segunda clave lleva `#2`).

## Fuera
- Extraer los últimos 4 dígitos de avisos y cartolas → parsers (ideas 008 y 010).
- Columna nueva en la base → fase 4 (migración de `charges`).
- Vista agrupada padre → líneas en el front → fases 7–9.
- Reglas del parser de cartola encontradas en la auditoría (monto en moneda de
  origen, excluir cuotas y pagos, separar el lugar pegado) → idea 010.
- Catálogo de servicios → idea 006 (en pausa hasta cerrar esta).

## Features afectadas
- `internal/shared/domain` (últimos 4 dígitos en `Charge` y `ParsedCharge`)
- `internal/subscriptions` (`domain` y los `usecases` que dependen de `GroupKey`)

## Dudas abiertas
- Ninguna. Decidido: modelo padre → líneas; separar por cuenta y últimos 4 dígitos;
  salto de monto en el mismo origen = cambio de plan; clave (b) con día típico.
- Cambios a `docs/ARCHITECTURE.md` al cerrar: formato de `GroupKey`, regla de líneas
  (reemplaza "umbrella → un grupo por monto"), últimos 4 dígitos en `Charge`, columna
  `card_last4` en el modelo de datos.
- Renumeración de la fase 3: avisos BCI → 008, extractor → 009, cartola BCI → 010.

## Historias
- [ ] [01 Últimos 4 dígitos de la tarjeta](01-card-last4.md)
- [ ] [02 Separar por origen](02-split-by-origin.md)
- [ ] [03 Encadenar líneas](03-chain-lines.md)
- [ ] [04 Clave estable de la línea](04-line-key.md)
- [ ] [05 `GroupCharges` produce líneas](05-group-charges-lines.md)
- [ ] [06 Resumen por línea con los escenarios reales](06-summary-by-line.md)
- [ ] [07 Alias de línea](07-line-alias.md)
