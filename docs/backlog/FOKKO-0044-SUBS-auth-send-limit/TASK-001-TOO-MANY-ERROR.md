# FOKKO-0044-SUBS/TASK-001 · Error de límite excedido
Feature: `internal/shared/domain`

Como usuario quiero un aviso claro si pedí demasiados códigos, para saber que debo esperar.

## Criterios
- Dada una llave, cuando se crea el error de límite, entonces envuelve `ErrTooMany` y lleva la llave y sus parámetros.
- Dado ese error, cuando se compara con `errors.Is(err, ErrTooMany)`, entonces es verdadero y no coincide con los demás.

## Fuera
- Nada.

## Notas
- Ninguna.
