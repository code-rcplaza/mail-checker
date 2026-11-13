# FOKKO-0094-SUBS/TASK-002B · Mostrar los errores del servidor
Feature: `frontend/src/features/auth/application`



## Criterios
- Dado un 409 `auth.email_taken`, cuando se envía, entonces muestra el mensaje en el campo del correo.
- Dado un 400 con `fields`, cuando se envía, entonces cada error aparece en su campo.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
