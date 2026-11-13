# FOKKO-0082D-SUBS/TASK-001 · Arranque que falla rápido
Feature: `cmd/api`

Como usuario quiero enterarme al arrancar si falta algo, para no descubrirlo en el primer uso.

## Criterios
- Dada configuración faltante o inválida, cuando arranca, entonces termina con un error explícito y código distinto de cero.
- Dada una migración que falla, cuando arranca, entonces termina con ese error y no sirve peticiones.
- Dado un arranque correcto, cuando está listo, entonces escribe un log de inicio sin secretos.

## Fuera
- Nada.

## Notas
- Ninguna.
