# FOKKO-0101-SUBS/TASK-003 · Conectar el hook con los componentes
Feature: `frontend/src/features/statements/presentation/screens`



## Criterios
- Dado el hook con un puerto doble, cuando se monta el contenedor, entonces pasa solo props a los componentes.
- Dado un componente de la pantalla, entonces no importa `application` ni `infrastructure`.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Contenedor en `presentation/screens`: llama hooks de `application` y pasa props.
