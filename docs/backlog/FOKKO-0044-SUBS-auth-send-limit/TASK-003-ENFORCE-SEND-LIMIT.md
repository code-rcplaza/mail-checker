# FOKKO-0044-SUBS/TASK-003 · Aplicar el límite al emitir
Feature: `internal/auth/usecases`

Como usuario quiero recibir un aviso en vez de más correos si me paso, para saber cuándo reintentar.

## Criterios
- Dados 3 envíos recientes al mismo correo, cuando se pide otro código, entonces devuelve `ErrTooMany` con `auth.too_many_sends` y no se envía nada.
- Dados 3 envíos a un correo y una petición para otro, cuando se pide el código, entonces no se limita (es por correo).
- Dado un envío permitido, cuando se emite, entonces queda registrado para contar.

## Fuera
- Nada.

## Notas
- Ninguna.
