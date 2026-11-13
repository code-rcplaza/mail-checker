# FOKKO-0096-SUBS/TASK-002 · Reenviar el código
Feature: `frontend/src/features/auth/application`

Como usuario quiero pedir otro código si no me llegó, para poder entrar.

## Criterios
- Dado el temporizador cumplido, cuando se toca "Reenviar código", entonces llama a `POST /api/auth/code/resend` y reinicia la cuenta.
- Dado un 429 `auth.too_many_sends`, cuando se pide, entonces muestra el mensaje y no reinicia la cuenta.
- Dado el reenvío correcto, cuando termina, entonces avisa con un `Toast`.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
