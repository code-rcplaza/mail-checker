# FOKKO-0098-SUBS/TASK-002 · Guardas de ruta
Feature: `frontend/src/features/auth/presentation/screens`

Como usuario quiero que mis datos no se vean sin sesión, para proteger mi privacidad.

## Criterios
- Dada una ruta protegida sin sesión, cuando se abre, entonces redirige a Login.
- Dadas Login o Registro con sesión, cuando se abren, entonces redirigen al Resumen.
- Dado un 401 en cualquier llamada con sesión aparente, cuando ocurre, entonces vuelve a Login y borra los datos cacheados.

## Fuera
- Nada.

## Notas
- Contenedor en `presentation/screens`: llama hooks de `application` y pasa props.
