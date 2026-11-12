# 006-01 · `Normalize`
Feature: `internal/subscriptions/adapters/catalog`

Como usuario quiero que los cobros de un mismo comercio se agrupen aunque el banco
escriba distinto su nombre (mayúsculas, puntos, prefijo del procesador de pago).

## Criterios
- Dado `"C. Verde Mall"`, cuando se normaliza, entonces `"C VERDE MALL"`.
- Dado `"Café Ñuñoa"`, cuando se normaliza, entonces `"CAFE NUNOA"` (sin acentos ni ñ).
- Dados `"DL*GOOGLE YOUTUBE"`, `"PAYU *UBER TRIP"`, `"MP *MERCADO LIB"`,
  `"MERCADOPAGO*ISTEC"` y `"DP *FALABELLA.C"`, cuando se normalizan, entonces quedan
  `"GOOGLE YOUTUBE"`, `"UBER TRIP"`, `"MERCADO LIB"`, `"ISTEC"` y `"FALABELLA C"`.
- Dados `"PedidosYa*Plus"`, `"ANTHROPIC* CLAUDE SUB"` y `"EA *ELECTRONIC ARTS"`, cuando
  se normalizan, entonces conservan el comercio: `"PEDIDOSYA PLUS"`,
  `"ANTHROPIC CLAUDE SUB"` y `"EA ELECTRONIC ARTS"`.
- Dado `"APPLE.COM/BILL   866-712-7753"`, cuando se normaliza, entonces
  `"APPLE COM BILL 866 712 7753"` (puntuación a espacio, espacios colapsados).
- Dado `"  dl*google  "`, cuando se normaliza, entonces `"GOOGLE"` (el prefijo se
  reconoce sin importar mayúsculas ni espacios alrededor).

## Fuera
- `Identify` → 006-02. `LooksLikeSubscription` → 006-03.

## Notas
- Descriptores tomados de cartolas BCI reales (nombres de empresas, sin datos
  personales).
- Prefijos como lista cerrada con nombre en el paquete; agregar uno es un commit.
