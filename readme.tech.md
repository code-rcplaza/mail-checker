# Fokko Subs · Documentación técnica

Punto de entrada para desarrolladores. El diseño vigente (capas, puertos, flujos y modelo de datos)
está en [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md), que es la fuente de verdad. Lo que entra en
la primera versión está en [`docs/mvp.md`](docs/mvp.md). Este archivo no repite esos documentos:
dice cómo está armado el código y cómo se trabaja con él.

## Stack

| Parte | Tecnología |
|---|---|
| Backend | Go (módulo `github.com/code-rcplaza/mail-checker/backend`, Go 1.27) |
| Router HTTP | `github.com/go-chi/chi/v5` |
| Base de datos | SQLite (`modernc.org/sqlite`, sin cgo), migraciones con `github.com/pressly/goose/v3` |
| Lectura de PDF | `github.com/ledongthuc/pdf` |
| Frontend | React + TypeScript, Tailwind v4, TanStack Router y Query, react-hook-form + zod (pendiente) |
| Despliegue | Docker (pendiente) |

Dependencias aprobadas por el equipo. Agregar una nueva requiere preguntar antes.

## Estructura del backend

```
backend/
├── cmd/api/            Punto de entrada: config, arranque y apagado limpio
└── internal/
    ├── shared/         Dinero, monedas, cargos, errores con llave i18n, SQLite, HTTP, reloj
    ├── accounts/       Cuentas de correo y sus credenciales cifradas (post-MVP)
    ├── emailsync/      Sincronización de correos de compra (post-MVP)
    ├── statements/     Cartolas: extracción de texto, parser de BCI, importación
    └── subscriptions/  Detección de suscripciones, resumen, decisiones y ajustes por usuario
```

Cada feature sigue la misma forma:

```
internal/<feature>/
├── domain/      Reglas puras; sin infraestructura
├── ports/       Interfaces que la feature necesita del exterior
├── usecases/    Casos de uso: orquestan puertos y dominio
└── adapters/    Implementaciones concretas (SQLite, PDF, parsers de banco, HTTP)
```

Reglas de dependencia:

- `domain` y `usecases` no importan `adapters` ni paquetes de infraestructura.
- Una feature no importa a otra; si necesita datos de otra, lo hace desde un adaptador puente
  (por ejemplo `emailsync/adapters/accounts`).
- `cmd/api` es la raíz de composición: arma las dependencias y no contiene lógica.

## Convenciones

- **Usuario en cada operación.** Todo dato es de un usuario. Los puertos de lectura y escritura
  reciben `userID` justo después de `ctx`, y cada consulta filtra por él.
- **Montos en centésimas**, como `shared.Money` (int64). Nunca floats. CLP 5.990 = `599000`.
- **Fechas:** unix en segundos, UTC, en la base y en la API. `time.Time` dentro de Go; la
  conversión vive solo en los adaptadores.
- **El tiempo se inyecta** (`now func() time.Time`). Nada de `time.Now()` en `domain` ni `usecases`.
- **Errores con contexto:** `fmt.Errorf("accion: %w", err)`. Nunca `panic` en flujo normal, nunca
  un error tragado.
- **Mensajes al usuario:** el backend devuelve llaves i18n (`feature.motivo`), no texto. Las
  llaves son constantes en `domain`. Los logs sí van en texto, con `slog`.
- **Sin datos personales en logs, fixtures ni commits.** Los PDFs reales no entran al repositorio;
  los fixtures son texto anonimizado escrito a mano.
- **Inmutabilidad:** las funciones no modifican slices ni mapas recibidos; devuelven copias.
- **Validación en el borde:** los constructores de `domain` (`NewCharge`, `NewService`,
  `NewManualSubscription`, `NewLineOverride`) validan y devuelven error. No existen valores de
  dominio inválidos.
- **Configuración:** variables de entorno leídas una vez al arrancar. Un valor inválido detiene la
  app con un error explícito.

## Configuración

| Variable | Por defecto | Qué hace |
|---|---|---|
| `MAILCHECKER_PORT` | `8080` | Puerto HTTP, entre 1 y 65535 |
| `MAILCHECKER_DB_PATH` | `data/mail-checker.db` | Archivo de SQLite; no puede ser vacío |

## Base de datos

- SQLite con `PRAGMA journal_mode=WAL`, `foreign_keys=ON`, `busy_timeout` y `synchronous=FULL`.
  Un solo escritor.
- Cada migración de goose corre en una transacción. La migración base (`00001_baseline.sql`) no
  crea tablas; cada feature agrega las suyas desde la `00002`.
- Las transacciones viven dentro de los métodos de los repositorios. Los casos de uso no las manejan.
- Dentro de una transacción solo se usa `tx`, nunca el repositorio ni `db`: la base tiene una sola
  conexión y usar `db` ahí se bloquea.
- Las escrituras son idempotentes (`UNIQUE` + `INSERT … ON CONFLICT DO NOTHING`): importar o
  sincronizar dos veces no duplica datos.
- Cada tabla con datos de usuario lleva `user_id` con `ON DELETE CASCADE`.

## Lectura de cartolas

1. `PDFExtractor` lee el texto del PDF. Las letras de los bancos vienen una por una y sin ancho,
   así que las palabras y columnas se reconstruyen por posición. Un PDF dañado devuelve error; un
   `panic` de la librería se convierte en error.
2. Cada parser de banco reconoce su cartola (`Detect`), devuelve sus compras con el número de
   línea (`Parse`), la fecha de facturación (`BillingDate`) y, si existe, el par de traspasos
   (`TransferPair`) que da el tipo de cambio.
3. `ImportStatement` elige el banco (el indicado o el único que lo reconoce), toma la fecha del
   archivo o la del usuario, guarda el tipo de cambio derivado y la cartola con sus cargos.

La cartola de BCI tiene dos estados: el nacional (compras en pesos) y el internacional (origen y
US$). Un cobro en pesos dentro del internacional se reconoce por la razón entre ambos montos.

## Pruebas

- **Dominio y casos de uso:** tests de tabla con dobles en memoria. Una prueba por criterio.
- **Parsers:** textos anonimizados en `testdata/`, con una prueba de tabla por caso.
- **Repositorios:** contra una base SQLite temporal real, con `-timeout 60s`: un bloqueo debe
  fallar en un minuto, no en diez.
- Solo la biblioteca estándar `testing`, sin testify.
- Un bug se reproduce con un test antes de corregirlo.

Comandos, desde `backend/`:

```sh
go build ./...
go vet ./...
staticcheck ./...        # instalado en $(go env GOPATH)/bin
gofmt -l .               # no debe listar nada
go test -count=1 -timeout 60s ./...
```

Algo está listo solo cuando build, vet, staticcheck, gofmt y los tests pasan con la salida a la vista.

## Flujo de trabajo

El backlog vive en `docs/backlog/`: cada idea es `FOKKO-XXXX-SUBS-slug/` con su `idea.md` y una
tarea por archivo `TASK-XXX-NOMBRE.md`. El orden es `/idea` → `/stories` → `/story`, con TDD: plan,
rojo, verde, refactor y verificación. La rama de trabajo es `dev`; `main` es producción y nunca
recibe commits directos.

Los commits van en inglés e identifican la tarea (`FOKKO-XXXX-SUBS/TASK-XXX`) en el cuerpo.

## Estado

Implementado en el dominio y los casos de uso del MVP: catálogo de servicios, detección,
cartolas de BCI, tipo de cambio, resumen, manuales, alzas, ajustes por línea y decisiones del
usuario.

Pendiente: la capa HTTP de cada caso de uso, la persistencia SQLite de las features nuevas, el
frontend y el empaquetado en Docker. El avance exacto está en `docs/backlog/README.md`.
