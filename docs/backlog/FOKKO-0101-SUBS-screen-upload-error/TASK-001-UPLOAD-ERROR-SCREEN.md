# FOKKO-0101-SUBS/TASK-001 · Pantalla de error
Feature: `frontend/src/features/statements/presentation/components`

Como usuario quiero entender por qué falló y qué hacer, para no quedarme trabado.

## Criterios
- Dada una subida fallida, cuando se renderiza, entonces muestra el mensaje traducido de la llave y una acción principal.
- Dado "Intentar de nuevo", cuando se toca, entonces vuelve a Elegir fuente.
- Dado el error, cuando se muestra, entonces no aparece ningún texto del archivo del usuario.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
