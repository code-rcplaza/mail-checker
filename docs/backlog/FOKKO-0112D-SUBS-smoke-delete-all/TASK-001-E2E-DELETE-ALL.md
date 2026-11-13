# FOKKO-0112D-SUBS/TASK-001 · Eliminar todo
Feature: `cmd/api`

Como usuario quiero que eliminar mi cuenta borre todo, para no dejar rastro.

## Criterios
- Dado un usuario con cartolas, decisiones, alias y manuales, cuando elimina su cuenta, entonces no queda ninguna fila suya en la base.
- Dada la sesión vieja, cuando se usa después, entonces responde 401.
- Dado otro usuario, cuando el primero elimina su cuenta, entonces sus datos siguen intactos.

## Fuera
- Nada.

## Notas
- Ninguna.
