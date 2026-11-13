# FOKKO-0095-SUBS/TASK-001B · Mostrar los errores del servidor
Feature: `frontend/src/features/auth/application`



## Criterios
- Dado un 404 `auth.user_not_found`, cuando se envía, entonces muestra el mensaje con un enlace a Registro.
- Dado un 429 `auth.too_many_sends`, cuando se envía, entonces muestra el mensaje y no navega.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
