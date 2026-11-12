# 007-02 · Separar por origen
Feature: `internal/subscriptions/domain`

Como usuario quiero que las suscripciones pagadas desde otra cuenta de correo u otra
tarjeta nunca se mezclen con las mías.

## Criterios
- Dados cobros de un mismo padre de las cuentas A y B (misma tarjeta), cuando se
  separan por origen, entonces quedan dos orígenes: uno por cuenta.
- Dados cobros de la misma cuenta con tarjetas `4821` y `1234`, cuando se separan,
  entonces quedan dos orígenes: uno por tarjeta.
- Dado un cobro de cartola sin cuenta (sin aviso) que encadena (intervalo válido,
  monto ±25 %) con los cobros de la cuenta A, cuando se separan, entonces se une al
  origen de A.
- Dado un cobro de cartola sin cuenta que no encadena con ningún origen, cuando se
  separan, entonces queda en un origen propio.
- Dado un cobro con `CardLast4` vacío, cuando se separa, entonces la tarjeta no
  distingue: se compara solo por cuenta.
- Dados un cobro de cartola y un aviso del mismo cobro con distintos últimos 4 dígitos,
  cuando se cruzan (`ChargesMatch`), entonces no calzan; si alguno no trae dígitos,
  la tarjeta no se compara.
- La función no modifica el slice recibido.

## Fuera
- Líneas dentro de un origen → 007-03.

## Notas
- Función pura; los cobros de prueba son sintéticos con montos de las cartolas reales.
