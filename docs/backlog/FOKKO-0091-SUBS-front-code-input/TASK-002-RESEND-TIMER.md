# FOKKO-0091-SUBS/TASK-002 · `ResendTimer`
Feature: `frontend/src/components/ui`

Como usuario quiero saber cuándo puedo pedir otro código, para no insistir.

## Criterios
- Dado el temporizador de 30 segundos, cuando corre, entonces muestra la cuenta regresiva y el botón está deshabilitado.
- Dado el tiempo cumplido, cuando termina, entonces el botón se habilita.
- Dado un reenvío, cuando se hace clic, entonces se reinicia la cuenta.

## Fuera
- Nada.

## Notas
- Los 30 segundos son del front; el límite real (3 por hora) lo aplica el back.
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
