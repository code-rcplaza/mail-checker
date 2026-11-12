# 006-02 · `Identify` y la tabla del catálogo
Feature: `internal/subscriptions/adapters/catalog`

Como usuario quiero que la app reconozca mis servicios conocidos sin confundirlos con
otros comercios, porque tomo decisiones con esa información.

## Criterios
- Dado el descriptor real de cada servicio del catálogo, cuando se llama `Identify`,
  entonces devuelve ese servicio:
  - `"NETFLIX.COM 866-579-7172"` y `"NETFLIX"` → Netflix;
  - `"Disney Plus Burbank"` → Disney+;
  - `"Spotify P1234ABCD Stockholm"` → Spotify;
  - `"ANTHROPIC* CLAUDE SUB"` → Claude;
  - `"HOSTGATOR LAS CONDES"` → HostGator;
  - `"PLAYIT.GG* PLAYIT.GG/C"` → playit.gg;
  - `"PedidosYa*Plus CL"` → PedidosYa Plus.
- Dados `"APPLE.COM/BILL 866-712-7753"`, `"DL*GOOGLE YOUTUBE"` y
  `"MICROSOFT. SUBSCRIPTION"`, cuando se llama `Identify`, entonces devuelve App Store,
  Google y Microsoft con `IsUmbrella = true`; los demás servicios tienen
  `IsUmbrella = false`.
- Dados `"NETFLIXER"`, `"MY NETFLIX"`, `"PedidosYa*KFC Vicun"`, `"PedidosYa*Propina"`,
  `"EA *ELECTRONIC ARTS"`, `"SEGUROS CLC S.A."`, `"MY MINI FACTORY LTD"`,
  `"NANONOBLE PTE. LTD."` y `"PAYU *UBER TRIP"`, cuando se llama `Identify`, entonces
  devuelve `false` (no calza a mitad de palabra ni a mitad del descriptor).
- Dado `"ANTHROPIC* API"`, cuando se llama `Identify`, entonces `false` (el patrón de
  Claude exige `ANTHROPIC CLAUDE`).

## Fuera
- Agrupación con el dominio → 006-04.

## Notas
- Calza solo si el descriptor normalizado empieza con el patrón y termina en límite de
  palabra; sin coincidencias aproximadas.
- El código de Spotify del ejemplo es inventado (el real identifica la cuenta).
