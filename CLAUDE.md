# Mail Checker

Detecta suscripciones a partir de los avisos de compra en Gmail y de las cartolas de
la tarjeta. Multi-banco (BCI primero). Uso personal y local. **Repo público.**

La arquitectura está en `docs/ARCHITECTURE.md` y es la fuente de verdad: si algo
la contradice, se pregunta antes de cambiarla.

## Forma de trabajo
- Responder en español.
- No asumir: stack, librerías, estructura, nombres o alcance se proponen y se
  aprueban antes de implementar.
- Trabajar por fases (plan en `docs/ARCHITECTURE.md`). Al terminar cada fase,
  detenerse y esperar revisión. No escribir código si no se ha pedido.
- Cambios destructivos (borrar archivos, reescribir historial, borrar repos)
  siempre con confirmación explícita.

## Principios

### KISS
- Lo más simple que funcione.
- Sin interfaces con una sola implementación; solo los puertos definidos en
  `docs/ARCHITECTURE.md` (parsers por banco, repos, lector de correo).
- Sin generics, reflection ni librerías fuera de las aprobadas sin preguntar.
- Go: errores con contexto (`fmt.Errorf("...: %w", err)`); nunca `panic` en flujo normal.
- Front: los datos del servidor viven solo en TanStack Query; nada de estado global
  duplicado (sin Redux/Zustand).

### YAGNI
- Solo se construye lo que pide la fase actual. Nada "por si acaso": ni
  configuraciones, ni flags, ni endpoints, ni campos sin uso hoy.
- Un segundo banco, Docker, etc. se implementan cuando se pidan; la arquitectura
  ya los permite, el código no se adelanta.

### DRY
- Regla de tres: dos copias se toleran; a la tercera se extrae.
- Una sola fuente de verdad: catálogo descriptor→servicio; reglas de detección
  (días, ±25 %, ventanas) como constantes con nombre en `domain`; `Money` y monedas
  en `shared/domain`; SQL solo en queries de sqlc.
- Entre features se prefiere duplicar antes que acoplar (manda la regla de no
  importar otra feature).
- Front: tipos de formularios con `z.infer<schema>`, no escritos a mano.

### Fail fast
- Validar en el borde: handler HTTP en Go y zod en el front. Lo inválido se
  rechaza ahí con `ErrInvalid` / 400 y un mensaje claro; no se propaga hacia adentro.
- Los constructores de `domain` validan sus invariantes y devuelven error; no
  existen valores de dominio inválidos.
- Al arrancar, configuración faltante o rota (llave, base, migraciones) detiene la
  app con un error explícito, no en el primer uso.
- Nunca tragarse errores (`_ = err`, `catch {}` vacío). Se manejan o se devuelven.

### Procesos como pasos
- Funciones pequeñas, de una sola cosa cada una, que se componen para que el
  código se lea como una lista de instrucciones. Ej. (en inglés, por la regla de
  idioma): `unixToTime`, `validateNotFuture`, `stripProcessorPrefix`,
  `groupByMerchant`.
- Cada función pequeña es testeable por separado; la que las junta solo las llama
  en orden.
- Un caso de uso se lee como receta: su método público solo orquesta y cada línea
  es un paso con nombre (`parsear`, `deduplicar`, `guardar`). La lógica de cada paso
  va en su propia función.
- Los pasos siguen el orden y los nombres de los flujos numerados de
  `docs/ARCHITECTURE.md`.
- Guard clauses y retorno temprano; no más de un nivel de `if` anidado; sin `else`
  después de `return`.
- Ningún paso esconde efectos: el nombre dice si lee, calcula o escribe. Los
  cálculos son funciones puras en `domain`; la E/S queda en los puertos.
- Front: handlers y hooks llaman pasos con nombre; sin lógica inline en el JSX.

### Inmutabilidad
- **Front (estricta)**: `readonly` en tipos, `as const`, nunca mutar `data` de
  TanStack Query, props ni estado; siempre copias nuevas.
- **Back (pragmática)**: inmutable siempre que se pueda y no genere un problema.
  - Preferir valores y funciones que devuelven un valor nuevo (sobre todo en `domain`:
    `Money`, cálculos de detección).
  - No modificar slices ni mapas recibidos como argumento; copiar antes de ordenar o
    hacer `append`.
  - Punteros permitidos cuando son lo natural en Go (repos que asignan el ID,
    structs grandes, opcionales).
  - Sin estado global mutable; el registro de bancos se arma en `main.go` y después
    solo se lee.

### Corrección
- Bug → primero un test que lo reproduce, después el arreglo.
- Parsers y detección se escriben con test primero.
- Go: `ctx` primer parámetro; IMAP y HTTP siempre con timeout.
- Go: el tiempo se inyecta (`now func() time.Time`); nada de `time.Now()` en `domain`
  ni en `usecases`.
- TS: `strict` + `noUncheckedIndexedAccess`; sin `any`, sin `!`, sin `as` para tapar
  errores de tipo.

### Legibilidad
- Código (identificadores) y mensajes de commit en inglés; docs, comentarios de
  diseño y UI en español.
- Nombres descriptivos de lo que la cosa **es** o **hace**, no conceptos:
  - Funciones: verbo + objeto (`ParseBCIEmail`, `DetectMonthlyCharges`,
    `StripProcessorPrefix`).
  - Tipos y variables: el sustantivo de lo que son (`Charge`, `StatementLine`,
    `monthlyCharges`).
  - Booleanos como pregunta (`isActive`, `hasGap`).
  - Prohibidos: `Manager`, `Service`, `Processor`, `Helper`, `Handler` a secas,
    `utils`, `common`, `data`, `info`, `item`, `obj`.
  - Sin abreviaturas salvo las idiomáticas de Go (`ctx`, `err`, `id`, `db`, `tx`, `i`).
  - Si un nombre necesita comentario para entenderse, se cambia el nombre.
- Comentarios explican el porqué, no el qué.
- Función de más de ~40 líneas o archivo de más de ~300: revisar si se parte
  (señal, no límite duro).
- Sin `TODO` sueltos: se hace ahora o va a "Pendiente" en `docs/ARCHITECTURE.md`.

### Privacidad (repo público)
- Logs estructurados con `slog`. Nunca se loguean contraseñas, cuerpos de correo,
  direcciones reales ni la llave.
- Fixtures de test se anonimizan antes de crearlos, no después.

### Verificación
- Go: `go vet` + `staticcheck`. Front: `biome lint` + `tsc --noEmit`.
- Cuando exista código: evaluar un hook `Stop` que corra los tests.

### ACID (SQLite)
- Atomicidad: cada método de repo que escribe más de una fila abre y cierra su
  transacción dentro del adaptador. Los casos de uso no manejan transacciones.
- Una transacción por unidad de trabajo: si falla la sincronización de una cuenta,
  se revierte solo esa; las demás siguen.
- Idempotencia: resincronizar o reimportar no duplica (`UNIQUE` +
  `INSERT … ON CONFLICT DO NOTHING`).
- Consistencia también en la base: `PRAGMA foreign_keys=ON`, `NOT NULL`, `CHECK`
  en montos y monedas.
- Aislamiento: WAL, `busy_timeout`, un solo escritor.
- Durabilidad: `synchronous=FULL`.
- Cada migración de goose corre en una transacción.

## Loop de desarrollo
`/idea` → `/stories NNN` → `/story NNN-MM` (skills en `.claude/skills/`).
- Paradas para revisión: después de la idea, de las historias, del plan, de RED y
  de VERIFY.
- Cada fase de `docs/ARCHITECTURE.md` se trabaja como una idea.
- Referencias `NNN-MM` (idea 3 dígitos, historia 2) en historias y commits.

### Organización de los `.md`
- Solo en tres lugares: raíz (`README.md`, `CLAUDE.md`), `docs/ARCHITECTURE.md` y
  `docs/backlog/**`. Cualquier otro `.md` se pregunta antes (nada de `notes.md`,
  `plan.md`, `TODO.md`; lo temporal va al scratchpad o a Engram).
- `docs/backlog/NNN-slug/idea.md` + `MM-slug.md` por historia. Slugs en inglés,
  kebab-case; contenido en español.
- Estado en un solo lugar: ideas en `docs/backlog/README.md`; historias en el
  checklist de su `idea.md`. Los archivos de historia no llevan estado.
- Ideas terminadas no se mueven ni se borran.
- El backlog es historial; `docs/ARCHITECTURE.md` es la verdad vigente. Si una
  historia cambia el diseño, se lleva a `ARCHITECTURE.md` al cerrarla.

## Reglas para Claude (no arrancarse)
- **Antes de escribir código**: decir qué archivos se tocan y qué se hace en cada
  uno, y esperar el ok.
- **Alcance cerrado**: solo los archivos del plan; si hace falta otro, se pregunta.
  Nada "de paso" (refactors, renombres, formateo ajeno a la tarea).
- **Siempre se pregunta antes de**: agregar dependencias, editar
  `docs/ARCHITECTURE.md`, `CLAUDE.md` o `.claude/`, commitear, borrar archivos.
- **"Listo" significa**: build, vet, tests y lint pasan, con la salida mostrada.
  Nunca decir "listo" sin haberlo corrido.
- **Tests intocables**: nunca se modifica, salta ni borra un test para que pase. Se
  arregla el código o se explica por qué el test está mal.
- **Dos intentos**: si algo falla dos veces por lo mismo, detenerse y explicar en
  vez de seguir probando a ciegas.
- **No inventar datos reales**: si falta un ejemplo real (formato de un banco,
  cartola), se pide.
- **Diffs chicos**: si un cambio pasa de ~300 líneas, se parte en pasos.
- **Al terminar**: resumen de archivos cambiados, por qué, y qué quedó fuera.
- **Ante la duda entre dos interpretaciones, se pregunta.**

## Git
- `main` = producción. Se trabaja en `dev`. Nunca commitear ni pushear a `main`.
- Commit y push solo cuando se pida.
- Repo público: nunca commitear datos personales (correos reales, cartolas,
  contraseñas, `.env`). Los ejemplos para tests se anonimizan.

## Backend (Go)
- Features verticales: `internal/<feature>/{adapters,domain,ports,usecases}`.
- `domain` y `usecases` no importan infraestructura ni `adapters`.
- Una feature no importa otra salvo desde un adaptador puente.
- Paquetes en minúscula sin `_` ni `-`; archivos en snake_case.
- Montos siempre como `Money` (centésimas, int64), nunca float.
- Fechas: unix en segundos (UTC) en la base y en la API; `time.Time` dentro de Go.
  La conversión vive solo en los adaptadores sqlite y http.
- IDs: `id` autoincremental solo para uso interno (nunca sale por la API);
  `public_id` UUID v7 generado en Go para lo que se expone. En el JSON, `id` = `public_id`.
- Tests: casos de uso con dobles en memoria; parsers con textos reales anonimizados.
- `gofmt` corre solo después de cada edición (hook en `.claude/settings.json`).

## Frontend (React + TS)
- shadcn/ui, Tailwind v4, mobile first (base 375 px, `md:`/`lg:` para escritorio).
- TanStack Router + TanStack Query; formularios con react-hook-form + zod.
- Estructura por features, reflejando el backend.
- Formato con Biome (hook en `.claude/settings.json`, activo cuando exista `frontend/`).

## Comandos
(se completan cuando exista el código: test, build, dev, sqlc, migraciones)
