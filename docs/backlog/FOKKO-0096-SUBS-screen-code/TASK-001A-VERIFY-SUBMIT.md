# FOKKO-0096-SUBS/TASK-001A · Enviar los seis dígitos
Feature: `frontend/src/features/auth/application`



## Criterios
- Dados seis dígitos, cuando se completan, entonces llama a `POST /api/auth/code/verify` solo una vez.
- Dado un código correcto, cuando se verifica, entonces navega a Bienvenida si es el primer acceso, o al Resumen si no.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
