# FOKKO-0028-SUBS/TASK-001 · Puerto de decisiones de alza
Feature: `internal/subscriptions/ports`

Como usuario quiero que mi decisión sobre un alza se recuerde, para no ver siempre el mismo aviso.

## Criterios
- Dado el puerto, cuando se guarda un descarte con clave de línea y monto nuevo, entonces `List` lo devuelve para ese usuario.
- Dado el mismo descarte guardado dos veces, cuando se guarda, entonces no se duplica.

## Fuera
- Nada.

## Notas
- Ninguna.
