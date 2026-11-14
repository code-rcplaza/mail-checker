# FOKKO-0097-SUBS · Pantalla Bienvenida

## Problema
Bienvenida saluda a quien entra por primera vez y lo lleva a elegir fuente. No llama a ningún endpoint propio.

## Depende de
- FOKKO-0098-SUBS y FOKKO-0092-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Pantalla `WelcomeScreen`.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/auth/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/auth`)

## Dudas abiertas
- TASK-001C queda abierta: no hay navegador en la máquina para revisar a 390, 768 y 1280 px.
- Resuelta: el handoff habla de correos de cobro; el MVP lee cartolas, así que los textos dicen
  "estado de cuenta" y sale "Prefiero subir un estado de cuenta" (sería el mismo botón). El
  saludo "Hola, {nombre}" va sobre el título, que el diseño no trae con nombre.
- Resuelta: con sesión lleva avatar y navegación (como el escritorio del handoff), sin sección marcada.

## Tareas
- [x] [TASK-001A Saludar con el nombre](TASK-001A-WELCOME-GREETING.md)
- [x] [TASK-001B Botón hacia Elegir fuente](TASK-001B-WELCOME-CTA.md)
- [ ] [TASK-001C Revisión visual en tres anchos](TASK-001C-WELCOME-LAYOUT.md)
- [x] [TASK-002 Conectar el hook con los componentes](TASK-002-CONTAINER-WELCOME.md)
