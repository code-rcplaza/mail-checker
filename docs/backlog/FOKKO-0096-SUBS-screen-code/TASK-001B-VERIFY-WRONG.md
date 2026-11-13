# FOKKO-0096-SUBS/TASK-001B · Mostrar un código inválido
Feature: `frontend/src/features/auth/application`



## Criterios
- Dado un `auth.code_invalid`, cuando se verifica, entonces muestra el error y vacía los casilleros.
- Dado un `auth.code_expired` o `auth.code_attempts_exceeded`, cuando se verifica, entonces ofrece pedir otro código.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
