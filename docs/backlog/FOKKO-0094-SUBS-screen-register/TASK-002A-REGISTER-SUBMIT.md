# FOKKO-0094-SUBS/TASK-002A · Enviar el registro
Feature: `frontend/src/features/auth/application`



## Criterios
- Dados datos válidos, cuando se envía, entonces llama a `POST /api/auth/register` y navega a Código con el correo.
- Dado el envío en curso, cuando se hace clic otra vez, entonces el botón está bloqueado.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
