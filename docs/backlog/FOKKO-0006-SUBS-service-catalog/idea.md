# FOKKO-0006-SUBS · Catálogo de servicios (fase 3, parte 1)

## Problema
La detección de suscripciones ya agrupa y decide con un `ServiceCatalog`, pero solo
existe como doble en los tests. Sin el catálogo real, ningún cobro se reconoce como
servicio conocido, los comercios desconocidos no se agrupan bien (mayúsculas,
puntuación, prefijos de procesador) y un cobro único nunca queda como posible. Esta
idea implementa el adaptador real, sin depender de ejemplos de banco.

## Alcance
- `subscriptions/adapters/catalog`: implementación de `ServiceCatalog` como **tabla en
  Go** (slice de reglas), sin archivos externos.
  - `Normalize(merchant)`:
    1. quita el prefijo de procesador de pago, de una **lista cerrada** vista en cartolas
       reales: `PAYU *`, `DL*` (dLocal), `MP *` y `MERCADOPAGO*` (Mercado Pago), `DP *`.
       No se quita cualquier `X*`: en `PEDIDOSYA*PLUS`, `ANTHROPIC* CLAUDE` o
       `EA *ELECTRONIC ARTS` lo que va antes del `*` es el comercio;
    2. pasa a mayúsculas y quita acentos;
    3. reemplaza puntuación por espacio y colapsa espacios.
  - `Identify(merchant)`: el descriptor normalizado **empieza con** el patrón de una
    regla, en palabras completas (`NETFLIX` calza `NETFLIX COM`, no `NETFLIXER`); sin
    coincidencias aproximadas. Devuelve `Service{Name, IsUmbrella}`; umbrella para
    tiendas de apps.
  - **Criterio de precisión**: un falso positivo (llamar Netflix a algo que no lo es) es
    peor que no reconocerlo, porque un comercio desconocido igual se detecta como
    posible por su cadencia y el usuario lo confirma. Solo entran al catálogo
    descriptores verificados en cobros reales.
  - `LooksLikeSubscription(merchant)`: verdadero si el descriptor normalizado contiene
    como palabra completa `SUBSCRIPTION`, `MEMBERSHIP`, `SUSCRIPCION`, `MEMBRESIA`,
    `PREMIUM` o `PLAN`.
- **Catálogo inicial** (todos vistos en cartolas BCI reales):

  | Servicio | Patrón (normalizado) | Ejemplo real | Umbrella |
  |---|---|---|---|
  | Netflix | `NETFLIX` | `NETFLIX.COM 866-…`, `NETFLIX` | no |
  | Disney+ | `DISNEY PLUS` | `Disney Plus Burbank` | no |
  | Spotify | `SPOTIFY` | `Spotify P<código> Stockholm` | no |
  | Claude | `ANTHROPIC CLAUDE` | `ANTHROPIC* CLAUDE SUB` | no |
  | HostGator | `HOSTGATOR` | `HOSTGATOR` | no |
  | playit.gg | `PLAYIT GG` | `PLAYIT.GG* PLAYIT.GG/C` | no |
  | PedidosYa Plus | `PEDIDOSYA PLUS` | `PedidosYa*Plus` | no |
  | App Store | `APPLE COM BILL` | `APPLE.COM/BILL 866-…` | sí |
  | Google | `GOOGLE` | `DL*GOOGLE YOUTUBE` | sí |
  | Microsoft | `MICROSOFT` | `MICROSOFT. SUBSCRIPTION`, `MICROSOFT MICROSO` | sí |

  - Umbrella cuando bajo un mismo nombre hay varios cobros distintos (Google cobra
    dos montos mensuales distintos como `GOOGLE YOUTUBE`; Microsoft y Apple venden
    muchas cosas): cada monto es su propio grupo.
  - **Fuera del catálogo a propósito** (comercios que también venden cosas sueltas o
    no se sabe qué son): `EA *ELECTRONIC ARTS`, `SEGUROS CLC`, `MY MINI FACTORY`,
    `NANONOBLE`. Si cobran con cadencia, igual aparecen como posibles.
  - `PedidosYa*KFC …` y `PedidosYa*Propina` no calzan `PEDIDOSYA PLUS`: quedan como
    comercios desconocidos.
- Tests con estos descriptores (nombres de empresas, sin datos personales).

## Fuera
- Parser de avisos BCI → post-MVP (lectura de correo).
- Extractor de archivos (`ledongthuc/pdf`, `excelize`, CSV, texto) → FOKKO-0009-SUBS.
- Parser de cartola BCI → FOKKO-0011A-SUBS y FOKKO-0011B-SUBS y FOKKO-0012-SUBS.
- Editar el catálogo desde la app → YAGNI (cambiarlo es un commit).
- Quitar sufijos de lugar o sucursal (`MALL`, `0123`) → no: podría juntar comercios
  distintos.

- En pausa hasta cerrar la idea FOKKO-0007-SUBS (líneas de suscripción): el catálogo alimenta esa
  agrupación.

## Features afectadas
- `internal/subscriptions` (solo `adapters/catalog`)

## Dudas abiertas
- Ninguna. La lista de servicios y prefijos sale de dos cartolas BCI reales (nacional e
  internacional); en internet no hay una lista confiable (los descriptores cambian por
  país y procesador).
- Para el parser de cartola, FOKKO-0011A-SUBS y FOKKO-0011B-SUBS y FOKKO-0012-SUBS (no para esta idea): la descripción de la cartola trae pegado el lugar
  (`NETFLIX.COM 866-579-7172`, `HOSTGATOR LAS CONDES`) y el aviso por correo puede
  traer otro texto para el mismo cobro. Separarlo es trabajo del parser.
- La tarea TASK-004 se reescribió tras la idea FOKKO-0007-SUBS: verifica padre y número de líneas
  (`service:`, `umbrella:`, `merchant:`), no las claves antiguas con monto y moneda.
- Cambio aplicado a `ARCHITECTURE.md`: el ejemplo de `Normalize` pasa a
  `"C. Verde Mall" → "C VERDE MALL"; "DLOCAL*SPOTIFY" → "SPOTIFY"`.

## Tareas
- [x] [TASK-001 `Normalize`](TASK-001-NORMALIZE.md)
- [x] [TASK-002 `Identify` y la tabla del catálogo](TASK-002-IDENTIFY.md)
- [x] [TASK-003 `LooksLikeSubscription`](TASK-003-LOOKS-LIKE-SUBSCRIPTION.md)
- [x] [TASK-004 Catálogo con la agrupación real](TASK-004-CATALOG-GROUPING.md)
