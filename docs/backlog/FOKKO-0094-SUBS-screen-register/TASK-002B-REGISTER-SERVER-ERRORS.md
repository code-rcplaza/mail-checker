# FOKKO-0094-SUBS/TASK-002B · Mostrar los errores del servidor
Feature: `frontend/src/features/auth/application`



## Criterios
- Dado un 429 `auth.too_many_sends`, cuando se envía, entonces muestra el mensaje y no navega.
- Dado un 400 con `fields`, cuando se envía, entonces cada error aparece en su campo.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
